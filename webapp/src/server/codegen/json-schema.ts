import { readFile } from "node:fs/promises";
import path from "node:path";

import {
  compile as compileSchemaInterface,
  type Options as JS2TSOptions,
} from "json-schema-to-typescript";

import { type IJSONSchema } from "#lib/json-schema";
import { reduceFolder } from "#server/lib";
import { fromJSON } from "#lib/json";

export interface IJSONSchemaCollection extends Record<string, IJSONSchema> {}

const schemaFileEnd = ".schema.json";
const metaSchemaFileName = `meta${schemaFileEnd}`;
let collection: IJSONSchemaCollection;

export async function collectJSONSchemas(schemaFolder: string) {
  if (collection) {
    return collection;
  }

  collection = await reduceFolder<IJSONSchemaCollection>(
    schemaFolder,
    {},
    async (schemaCollection, filePath, entry) => {
      // for now exclude meta schema
      const isSchema =
        entry.isFile() &&
        entry.name !== metaSchemaFileName &&
        entry.name.endsWith(schemaFileEnd);

      if (!isSchema) {
        return schemaCollection;
      }

      const schemaPath = path.format(filePath);
      const schemaJSON = await readFile(schemaPath, { encoding: "utf-8" });
      const jsonSchema = fromJSON<IJSONSchema>(schemaJSON);

      if (!jsonSchema.$id || !jsonSchema.title) {
        throw new Error(
          `JSON schema file should have "$id" and "title" properties at the top level and the schema at "${filePath}" misses either of them`
        );
      }

      schemaCollection[jsonSchema.$id] = jsonSchema;

      return schemaCollection;
    }
  );

  return collection;
}

interface ISchemaInterface {
  inputSchema: IJSONSchema;
  code: string;
}

const parserOptions: JS2TSOptions["$refOptions"] = {
  resolve: {
    http: {
      async read(file) {
        const schemaCopy = transformSchema(collection[file.url]);
        return schemaCopy;
      },
    },
  },
};

export async function createInterfaceFromSchema(
  inputSchema: IJSONSchema,
  options?: Partial<Omit<JS2TSOptions, "$refOptions" | "bannerComment">>
): Promise<ISchemaInterface> {
  const schemaCopy = transformSchema(inputSchema);
  const interfaceCode = await compileSchemaInterface(
    // @ts-expect-error draft-04 type
    schemaCopy,
    schemaCopy.$id,
    { ...options, bannerComment: "", $refOptions: parserOptions }
  );

  const schemaInterface: ISchemaInterface = {
    inputSchema: schemaCopy,
    code: interfaceCode,
  };

  return schemaInterface;
}

function transformSchema(inputSchema: IJSONSchema): IJSONSchema {
  const modifedSchema = Object.assign({}, inputSchema);
  modifedSchema.title = `I${modifedSchema.title}`;

  return modifedSchema;
}
