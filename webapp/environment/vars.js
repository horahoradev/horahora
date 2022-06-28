export const API_ENDPOINT = process.env.NEXT_PUBLIC_API_ENDPOINT;

if (API_ENDPOINT == null) {
  throw new Error("No NEXT_PUBLIC_API_ENDPOINT provided!");
}
