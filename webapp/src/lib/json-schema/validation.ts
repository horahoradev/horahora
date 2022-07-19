import Ajv, { type ValidateFunction, type DefinedError } from "ajv";
import { type JSONSchema7 } from "json-schema";

import { multilineString } from "#lib/strings";

export interface IJSONSchema extends JSONSchema7 {
  async: false;
}

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

const ajv = createAJV();

export function createValidator<SchemaInterface>(schemaID: string) {
  const validate: ValidateFunction<SchemaInterface> | undefined =
    ajv.getSchema<SchemaInterface>(schemaID);

  if (!validate) {
    throw new Error(`JSON Schema with "$id" "${schemaID}" doesn't exist.`);
  }

  return (inputJSON: unknown): inputJSON is SchemaInterface => {
    const result = validate(inputJSON);

    if (!result) {
      // errors key is always an array when validation is failed
      const errors = [...(validate.errors! as DefinedError[])];
      throw new ValidationError(errors, schemaID);
    }

    return true;
  };
}

function createAJV() {
  return new Ajv({
    schemas: [],
  });
}
