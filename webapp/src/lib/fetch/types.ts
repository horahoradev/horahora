import { PUBLIC_API_URL } from "#environment/derived";

export const formHeader = ["content-type", "application/x-www-form-urlencoded"];

/**
 * URL constructor for public API endpoints.
 */
export class PublicAPIURL extends URL {
  /**
   *
   * @param path A pathname without search params which should start with a slash.
   * @param searchParams Search parameters of the URL.
   */
  constructor(path: string, searchParams?: URLSearchParams) {
    if (!path.startsWith("/")) {
      throw new Error("Public URL `path` argument should start with a slash.");
    }

    super(`${PUBLIC_API_URL.pathname}${path}`, PUBLIC_API_URL.origin);

    if (searchParams) {
      this.search = searchParams.toString();
    }
  }
}

export class FetchError extends Error {
  constructor(baseMessage: string, response: Response, body: any = []) {
    const message = [
      `${baseMessage}. Details:`,
      `Status: ${response.status}`,
      `Message: ${response.statusText}`,
      `Reason: ${body.message}`
    ].join("\n");
    super(message);
  }
}

export async function FetchErrorWithBody(baseMessage: string, response: Response) {
  let body: Uint8Array = await response.json();
  return new FetchError(baseMessage, response, body)
}
