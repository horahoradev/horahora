import { PublicAPIURL } from "#lib/fetch";

/**
 * @TODO proper args typing
 */
export async function uploadFile(formData: FormData) {
  const url = new PublicAPIURL("/upload");
  const response = await fetch(url, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    const message = [
      "Failed to upload the file. Details:",
      `Status: ${response.status}`,
      `Message: ${response.statusText}`,
    ].join("\n");
    throw new Error(message);
  }

  const newVideoID: number = await response.json();
  return newVideoID;
}
