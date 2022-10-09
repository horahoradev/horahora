import path from "node:path";
import { cwd } from "node:process";

import stringifyObject from "stringify-object";

import {
  collectJSONSchemas,
  type ICodegen,
  type ICodegenFunc,
} from "#server/codegen";
import { decapitalizeString, multilineString } from "#lib/strings";

const schemaFolder = ["..", "schema"];

const generateSchemaMap: ICodegenFunc = async () => {
  const schemaFolderPath = path.join(cwd(), ...schemaFolder);
  const schemaCollection = await collectJSONSchemas(schemaFolderPath);
  const results: string[] = [];
  const concreteExports: string[] = ["schemaMap"];
  const schemaMap: Record<string, string> = {};

  for await (const [schemaID, schema] of Object.entries(schemaCollection)) {
    const schemaName = `${decapitalizeString(schema.title!)}Schema`;
    const schemaObj = `export const ${schemaName} = ${stringifyObject(
      schema
    )} as const`;

    schemaMap[`[${schemaName}.$id]`] = schemaName;
    // add schema name to the index exports
    concreteExports.push(schemaName);
    results.push(schemaObj);
  }

  const schemaMapString = multilineString(
    ...Object.entries(schemaMap).map(([idRef, schemaName]) => {
      return `${idRef}: ${schemaName},`;
    })
  );
  results.push(`export const schemaMap = {${schemaMapString}} as const`);

  const codegen: ICodegen = {
    exports: {
      concrete: concreteExports,
    },
    result: multilineString(...results),
  };

  return codegen;
};

export default generateSchemaMap;
