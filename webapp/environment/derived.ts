import { PUBLIC_API_ORIGIN, PUBLIC_API_PREFIX } from "./vars";

export const PUBLIC_API_URL = new URL(PUBLIC_API_PREFIX, PUBLIC_API_ORIGIN);