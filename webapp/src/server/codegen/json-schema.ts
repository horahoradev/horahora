import { readFile } from "node:fs/promises";
import path from "node:path";

import { type IJSONSchema } from "#lib/json-schema";
import { reduceFolder } from "#server/lib";
import { fromJSON } from "#lib/json";

export interface IJSONSchemaCollection extends Record<string, IJSONSchema> {}

const schemaFileEnd = ".schema.json";
let collection: IJSONSchemaCollection | undefined = undefined;

export async function collectJSONSchemas(schemaFolder: string) {
  if (collection) {
    return collection;
  }

  collection = await reduceFolder<IJSONSchemaCollection>(
    schemaFolder,
    {},
    async (schemaCollection, filePath, entry) => {
      const isSchema = entry.isFile() && entry.name.endsWith(schemaFileEnd);

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
