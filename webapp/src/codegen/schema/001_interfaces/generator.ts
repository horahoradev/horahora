import path from "node:path";
import { cwd } from "node:process";

import {
  collectJSONSchemas,
  createInterfaceFromSchema,
  type ICodegen,
  type ICodegenFunc,
} from "#server/codegen";
import { multilineString } from "#lib/strings";

const schemaFolder = ["..", "schema"];

const generateJSONSchemas: ICodegenFunc = async () => {
  const schemaFolderPath = path.join(cwd(), ...schemaFolder);
  const schemaCollection = await collectJSONSchemas(schemaFolderPath);
  const results: string[] = [];
  const typeExports: string[] = [];

  for await (const [schemaID, schema] of Object.entries(schemaCollection)) {
    const schemaInterface = await createInterfaceFromSchema(schema, {
      declareExternallyReferenced: false,
    });

    // add interface name to the index exports
    typeExports.push(schemaInterface.inputSchema.title!);
    results.push(schemaInterface.code);
  }

  const codegen: ICodegen = {
    exports: {
      types: typeExports,
    },
    result: multilineString(...results),
  };

  return codegen;
};

export default generateJSONSchemas;
