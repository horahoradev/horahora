import { FetchError, formHeader, PublicAPIURL } from "#lib/fetch";
import {
  IArchivalRequest,
  type IArchivalEvent,
} from "#codegen/schema/001_interfaces";

export async function createNewArchivalRequest(formParams: URLSearchParams) {
  const url = new PublicAPIURL("/archiverequests");
  const headers = new Headers([formHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formParams,
  });

  if (!response.ok) {
    const error = NewAsyncFetchError(
      "Failed to create a new archival request",
      response
    );
    throw error;
  }

  const result: null = await response.json();

  return result;
}

interface IRequestInfo {
  ArchivalEvents: IArchivalEvent[];
  ArchivalRequests: IArchivalRequest[];
}

export async function getArchivalEvents(downloadID: any) {
  const url = new PublicAPIURL("/archiveevents/" + downloadID);
  const response = await fetch(url, {
    method: "GET",
  });

  if (!response.ok) {
    const error = new FetchError("Failed to fetch archival requests", response);
    throw error;
  }

  const result: IRequestInfo = await response.json();

  return result;
}

export async function getArchivalRequests() {
  const url = new PublicAPIURL("/archiverequests");
  const response = await fetch(url, {
    method: "GET",
  });

  if (!response.ok) {
    const error = new FetchError("Failed to fetch archival requests", response);
    throw error;
  }

  const result: IRequestInfo = await response.json();

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
    const error = new FetchError(
      "Failed to delete an archival request",
      response
    );
    throw error;
  }

  return null;
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
    const error = new FetchError(
      "Failed to retry an archival request",
      response
    );
    throw error;
  }

  const result: null = await response.json();

  return result;
}

export async function fetchDownloadsInProgress() {
  const url = new PublicAPIURL("/downloadsinprogress");
  const response = await fetch(url, {
    method: "GET",
  });

  if (!response.ok) {
    const error = new FetchError(
      "Failed to fetch downloads in progress",
      response
    );
    throw error;
  }

  const result = await response.json();

  return result;
}
