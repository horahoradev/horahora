import Ajv, { type ValidateFunction, type DefinedError } from "ajv";
import { type JSONSchema7 } from "json-schema";
import addFormats from "ajv-formats";

import { multilineString } from "#lib/strings";

export interface IJSONSchema extends JSONSchema7 {
  $id: string;
  title: string;
}

export interface ISchemaMap
  extends Record<Required<IJSONSchema>["$id"], IJSONSchema> {}

export class ValidationError extends Error {
  errors: DefinedError[];
  schemaID: string;

  constructor(errors: DefinedError[], schemaID: string) {
    const errorJSONs = errors!.map((error) => {
      return error.message!;
    });
    const message = multilineString(
      `Failed to validate according to schema "${schemaID}"`,
      ...errorJSONs
    );

    super(message);
    this.name = this.constructor.name;
    this.errors = errors;
    this.schemaID = schemaID;
  }
}

export function createAJV(schemaMap: ISchemaMap) {
  const ajv = new Ajv({
    schemas: Object.values(schemaMap),
  });

  addFormats(ajv);

  return ajv;
}

export function createValidator<SchemaInterface>(
  schema: IJSONSchema,
  ajv: ReturnType<typeof createAJV>
) {
  let validate: ValidateFunction<SchemaInterface> | undefined = undefined;

  return (inputJSON: unknown): inputJSON is SchemaInterface => {
    if (!validate) {
      validate = ajv.getSchema<SchemaInterface>(schema.$id);

      if (!validate) {
        throw new Error(
          `JSON Schema with "$id" "${schema.$id}" doesn't exist.`
        );
      }
    }

    const result = validate(inputJSON);

    if (!result) {
      // `errors` key is always an array when validation is failed
      const errors = [...(validate.errors! as DefinedError[])];
      throw new ValidationError(errors, schema.$id);
    }

    return true;
  };
}
