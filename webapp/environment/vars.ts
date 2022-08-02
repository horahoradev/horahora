/**
 * Available on client without any prefix.
 */
export const NODE_ENV = process.env.NODE_ENV;
export const PUBLIC_ORIGIN = process.env.NEXT_PUBLIC_ORIGIN!;
export const PUBLIC_API_ORIGIN = process.env.NEXT_PUBLIC_PUBLIC_API_ORIGIN!;
export const PUBLIC_API_PREFIX =
  process.env.NEXT_PUBLIC_PUBLIC_API_PATHNAME_PREFIX!;

export const PUBLIC_WEBSOCKET_ORIGIN =
  process.env.NEXT_PUBLIC_WEBSOCKET_ORIGIN!;
export const PUBLIC_WEBSOCKET_PATHNAME_PREFIX =
  process.env.NEXT_PUBLIC_WEBSOCKET_PATHNAME_PREFIX!;
