import { PUBLIC_API_URL } from "#environment/derived";

export const UserRank = {
  REGULAR: 0,
  TRUSTED: 1,
  ADMIN: 2,
} as const;

export type IUserRank = typeof UserRank[keyof typeof UserRank];

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
      throw new Error("Public URL argument should start with a slash.");
    }

    super(`${PUBLIC_API_URL.pathname}${path}`, PUBLIC_API_URL.origin);

    if (searchParams) {
      this.search = searchParams.toString();
    }
  }
}

export class FetchError extends Error {
  constructor(baseMessage: string, response: Response) {
    const message = [
      `${baseMessage}. Details:`,
      `Status: ${response.status}`,
      `Message: ${response.statusText}`,
    ].join("\n");
    super(message)
  }
}
