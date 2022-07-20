import Ajv, { type ValidateFunction, type DefinedError } from "ajv";
import { type JSONSchema7 } from "json-schema";

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
  return new Ajv({
    schemas: Object.values(schemaMap),
  });
}

export function createValidator<SchemaInterface>(
  schema: IJSONSchema,
  ajv: ReturnType<typeof createAJV>
) {
  const validate: ValidateFunction<SchemaInterface> | undefined =
    ajv.getSchema<SchemaInterface>(schema.$id);

  if (!validate) {
    throw new Error(`JSON Schema with "$id" "${schema.$id}" doesn't exist.`);
  }

  return (inputJSON: unknown): inputJSON is SchemaInterface => {
    const result = validate(inputJSON);

    if (!result) {
      // `errors` key is always an array when validation is failed
      const errors = [...(validate.errors! as DefinedError[])];
      throw new ValidationError(errors, schema.$id);
    }

    return true;
  };
}
