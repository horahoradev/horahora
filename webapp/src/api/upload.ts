import { PublicAPIURL } from "./types";

/**
 * @TODO proper args typing
 */
export async function uploadFile(formData: FormData) {
  const url = new PublicAPIURL("/upload");
  const response = await fetch(url, {
    method: "GET",
    body: formData,
  });

  if (!response.ok) {
    const message = [
      "Failed to upload the file. Error details:",
      `Status: ${response.status}`,
      `Message: ${response.statusText}`,
    ].join("\n");
    throw new Error(message);
  }

  const newVideoID: number = await response.json();
  return newVideoID;
}
