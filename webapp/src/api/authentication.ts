import { FetchError, formHeader, PublicAPIURL } from "./types";

/**
 * @TODO proper arg typing
 */
export async function registerAccount(formParams: URLSearchParams) {
  const url = new PublicAPIURL("/register");
  const headers = new Headers([formHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formParams,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to register an account", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}

export async function loginAccount(formParams: URLSearchParams) {
  const url = new PublicAPIURL("/login");
  const headers = new Headers([formHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formParams,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to log in", response);
    throw error;
  }

  const result = await response.json();

  return result;
}

export async function logoutAccount() {
  const url = new PublicAPIURL("/logout");
  const response = await fetch(url, {
    method: "POST",
  });

  if (!response.ok) {
    const error = new FetchError("Failed to log out", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}

export async function resetAccountPassword(formParams: URLSearchParams) {
  const url = new PublicAPIURL("/password-reset");
  const headers = new Headers([formHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formParams,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to reset password", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}
