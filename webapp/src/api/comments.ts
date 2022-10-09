import { FetchError, formHeader, PublicAPIURL } from "#lib/fetch";

export async function addComment(formParams: URLSearchParams) {
  const url = new PublicAPIURL("/comments/");
  const headers = new Headers([formHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formParams,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to add a comment", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}
