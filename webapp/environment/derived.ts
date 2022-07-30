import {
  PUBLIC_API_ORIGIN,
  PUBLIC_API_PREFIX,
  NODE_ENV,
  PUBLIC_ORIGIN,
} from "./vars";

export const IS_DEVELOPMENT = NODE_ENV === "development";
export const PUBLIC_URL = new URL(PUBLIC_ORIGIN);
export const PUBLIC_API_URL = new URL(PUBLIC_API_PREFIX, PUBLIC_API_ORIGIN);
