import { FetchError, multipartFormHeader, PublicAPIURL } from "./types";

/**
 * @TODO proper arg typing
 */
export async function registerAccount(formData: FormData) {
  const url = new PublicAPIURL("/register");
  const headers = new Headers([multipartFormHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formData,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to register an account", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}

export async function loginAccount(formData: FormData) {
  const url = new PublicAPIURL("/login");
  const headers = new Headers([multipartFormHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formData,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to login an account", response);
    throw error;
  }

  const result = await response.json();

  return result;
}

export async function logoutAccount() {
  const url = new PublicAPIURL("/logout");
  const response = await fetch(url, {
    method: "GET",
  });

  if (!response.ok) {
    const error = new FetchError("Failed to log out an account", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}

export async function resetAccountPassword(formData: FormData) {
  const url = new PublicAPIURL("/password-reset");
  const headers = new Headers([multipartFormHeader]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formData,
  });

  if (!response.ok) {
    const error = new FetchError("Failed to reset password", response);
    throw error;
  }

  const result: null = await response.json();

  return result;
}
