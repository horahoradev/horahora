import { PublicAPIURL } from "./types";

/**
 * @TODO proper arg typing
 */
export async function postRegister(formData: FormData) {
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
