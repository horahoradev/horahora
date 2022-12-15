import { PublicAPIURL, FetchError } from "./types";

import { logoutAccount } from "#lib/account";

export interface IAPIFetchArgs {
  pathname: ConstructorParameters<typeof PublicAPIURL>["0"];
  searchParams?: ConstructorParameters<typeof PublicAPIURL>["1"];
  baseErrorMessage?: string;
  options?: IAPIFetchOptions;
}

export interface IAPIFetchOptions extends Omit<RequestInit, "credentials"> {
  method?: "GET" | "POST";
}

export async function apiFetch<ResBody = never>({
  pathname,
  searchParams,
  baseErrorMessage = "Failed to fetch",
  options,
}: IAPIFetchArgs): Promise<ResBody> {
  const url = new PublicAPIURL(pathname, searchParams);
  const response = await fetch(url, { ...options, credentials: "include" });

  // Why the hell does the status say 200 if it's a 301??
  if (response.redirected) {
    console.log(`Redirecting to new location: ${response.url}`);
    location.replace(response.url);
  }

  if (!response.ok) {
    // @TODO: 403 status handling
    switch (response.status) {
      case 403: {
        await logoutAccount();
        // @TODO more graceful way of handling logout
        location.reload();
      }
      default: {
        const error = new FetchError(baseErrorMessage, response);
        throw error;
      }
    }
  }

  const data: ResBody = await response.json();

  return data;
}
