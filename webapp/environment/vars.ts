/**
 * Available on client without any prefix.
 */
export const NODE_ENV = process.env.NODE_ENV;
export const PUBLIC_API_ORIGIN = process.env.NEXT_PUBLIC_PUBLIC_API_ORIGIN!;
export const PUBLIC_API_PREFIX =
  process.env.NEXT_PUBLIC_PUBLIC_API_PATHNAME_PREFIX!;
export const API_ENDPOINT = process.env.NEXT_PUBLIC_API_ENDPOINT!;

if (API_ENDPOINT == null) {
  throw new Error("No NEXT_PUBLIC_API_ENDPOINT provided!");
}
