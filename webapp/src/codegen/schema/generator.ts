import path from "node:path";
import { cwd } from "node:process";

import stringifyObject from "stringify-object";

import {
  collectJSONSchemas,
  createInterfaceFromSchema,
  type ICodegen,
  type ICodegenFunc,
} from "#server/codegen";
import { decapitalizeString, multilineString } from "#lib/strings";

const schemaFolder = ["..", "schema"];

/**
 * @TODO validator gen
 */
const generateJSONSchemas: ICodegenFunc = async () => {
  const schemaFolderPath = path.join(cwd(), ...schemaFolder);
  const schemaCollection = await collectJSONSchemas(schemaFolderPath);
  const results: string[] = [];
  const typeExports: string[] = [];
  const concreteExports: string[] = [];

  for await (const [schemaID, schema] of Object.entries(schemaCollection)) {
    const schemaInterface = await createInterfaceFromSchema(schema, {
      declareExternallyReferenced: false,
    });

    const schemaName = `${decapitalizeString(schema.title!)}Schema`;
    const schemaObj = `export const ${schemaName} = ${stringifyObject(
      schema
    )} as const`;

    // add interface name to the index exports
    typeExports.push(schemaInterface.inputSchema.title!);
    // add schema name to the index exports
    concreteExports.push(schemaName);
    results.push(schemaInterface.code, schemaObj);
  }

  const codegen: ICodegen = {
    exports: {
      concrete: concreteExports,
      types: typeExports,
    },
    result: multilineString(...results),
  };

  return codegen;
};

export default generateJSONSchemas;
