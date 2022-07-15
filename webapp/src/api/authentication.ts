import { PublicAPIURL } from "./types";

/**
 * @TODO proper arg typing
 */
export async function registerAccount(formData: FormData) {
  const url = new PublicAPIURL("/register");
  const headers = new Headers([["content-type", "multipart/form-data"]]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formData,
  });

  if (!response.ok) {
    const message = [
      "Failed to register an account. Details:",
      `Status: ${response.status}`,
      `Message: ${response.statusText}`,
    ].join("\n");
    throw new Error(message);
  }

  const result: null = await response.json();
  return result;
}

export async function loginAccount(formData: FormData) {
  const url = new PublicAPIURL("/login");
  const headers = new Headers([["content-type", "multipart/form-data"]]);
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: formData,
  });

  if (!response.ok) {
    const message = [
      "Failed to login an account. Details:",
      `Status: ${response.status}`,
      `Message: ${response.statusText}`,
    ].join("\n");
    throw new Error(message);
  }

  const result = await response.json();
  return result;
}
