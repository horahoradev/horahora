import { FetchError, PublicAPIURL } from "./types";

import { type IAccountClient } from "#codegen/schema/001_interfaces";

export async function fetchAccountInfo() {
  const url = new PublicAPIURL("/currentuserprofile");
  const response = await fetch(url, {
    method: "GET",
  });

  if (!response.ok) {
    const error = new FetchError("Failed fetch account info", response);
    throw error;
  }

  const result: IAccountClient = await response.json();

  return result;
}
