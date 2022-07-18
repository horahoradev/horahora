import { readFile } from "node:fs/promises";
import path from "node:path";
import { cwd } from "node:process";

import { compile } from "json-schema-to-typescript";
import { type JSONSchema7 } from "json-schema";

import { handleError, isError } from "#lib/errors";
import { multilineString } from "#lib/strings";

interface IMetaSchema extends JSONSchema7 {
  $id: string;
  title: string;
}

const metaSchemaFilename = "meta.schema.json";

async function generateMetaSchema() {
  try {
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
        format: false,
        declareExternallyReferenced: false,
      }
    );

    return metaSchemaInterface;
  } catch (error) {
    handleError(error, "Failed to generate meta schema interface");
  }
}

function transformSchema(inputSchema: IMetaSchema): IMetaSchema {
  const modifedSchema = Object.assign({}, inputSchema);
  modifedSchema.title = `I${modifedSchema.title}`;

  return modifedSchema;
}

export default generateMetaSchema;
