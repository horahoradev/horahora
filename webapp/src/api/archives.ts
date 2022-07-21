import { FetchError, formHeader, PublicAPIURL } from "./types";

export async function createNewArchivalRequest(formParams: URLSearchParams) {
  const url = new PublicAPIURL("/archiverequests");
  const headers = new Headers([formHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formParams,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to create a new archival request", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}

export async function deleteArchivalRequest(formParams: URLSearchParams) {
  const url = new PublicAPIURL("/delete-archiverequest");
  const headers = new Headers([formHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formParams,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to delete an archival request", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}

export async function retryArchivalRequest(formParams: URLSearchParams) {
  const url = new PublicAPIURL("/retry-archiverequest");
  const headers = new Headers([formHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formParams,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to retry an archival request", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}
