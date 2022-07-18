import { multilineString } from "#lib/strings";

export function handleError(error: unknown, errorMessage: string): void {
  if (!isError(error)) {
    throw error;
  }

  const message = multilineString(
    errorMessage,
    `Reason: ${error.message}`
  );
  throw new Error(message, { cause: error });
}

export function isError(error: unknown): error is Error {
  return error instanceof Error
}
