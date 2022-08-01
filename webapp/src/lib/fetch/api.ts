import { PublicAPIURL, FetchError } from "./types";

export interface IAPIFetchArgs {
  pathname: ConstructorParameters<typeof PublicAPIURL>["0"];
  searchParams?: ConstructorParameters<typeof PublicAPIURL>["1"];
}

export interface IAPIFetchOptions extends RequestInit {
  baseErrorMessage?: string;
}

export async function apiFetch<ResBody = never>(
  { pathname, searchParams }: IAPIFetchArgs,
  { baseErrorMessage = "Failed to fetch", ...fetchOptions }: IAPIFetchOptions
): Promise<ResBody> {
  const url = new PublicAPIURL(pathname, searchParams);
  const response = await fetch(url, fetchOptions);

  if (!response.ok) {
    // @TODO: 403 status handling
    switch (response.status) {
      default: {
        const error = new FetchError(baseErrorMessage, response);
        throw error;
      }
    }
  }

  const data: ResBody = await response.json();

  return data;
}
