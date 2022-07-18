import { readFile } from "node:fs/promises";
import path from "node:path";
import { cwd } from "node:process";

import { compile } from "json-schema-to-typescript";
import { type JSONSchema7 } from "json-schema";

import { type ICodegen, type ICodegenFunc } from "#server/codegen";

interface IMetaSchema extends JSONSchema7 {
  $id: string;
  title: string;
}

const metaSchemaFilename = "meta.schema.json";

const generateMetaSchema: ICodegenFunc = async () => {
  const metaSchemaPath = path.join(cwd(), "..", "schema", metaSchemaFilename);
  const metaSchemaJSON = await readFile(metaSchemaPath, {
    encoding: "utf-8",
  });
  const metaSchema: IMetaSchema = JSON.parse(metaSchemaJSON);
  const modifedSchema = transformSchema(metaSchema);

  const metaSchemaInterface = await compile(
    // @ts-expect-error JSONSchema4 incompatibility
    modifedSchema,
    modifedSchema.$id,
    {
      additionalProperties: false,
      bannerComment: "",
      format: false
    }
  );

  const typeExports = [modifedSchema.title];

  const codegen: ICodegen = {
    exports: {
      types: typeExports,
    },
    result: metaSchemaInterface,
  };

  return codegen;
};

function transformSchema(inputSchema: IMetaSchema): IMetaSchema {
  const modifedSchema = Object.assign({}, inputSchema);
  modifedSchema.title = `I${modifedSchema.title}`;

  return modifedSchema;
}

export default generateMetaSchema;
