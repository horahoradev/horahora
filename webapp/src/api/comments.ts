import { FetchError, PublicAPIURL } from "./types";

export async function addComment(formData: FormData) {
  const url = new PublicAPIURL("/comments");
  const headers = new Headers([["content-type", "multipart/form-data"]]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formData,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to add a comment", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}
