import Cookies from "js-cookie";
import type { CookieAttributes } from "js-cookie";

const cookieKeys = {
  THEME: "theme",
} as const;

/**
 * Only includes client-accessible cookies.
 */
type ICookieKey = typeof cookieKeys[keyof typeof cookieKeys];

const defaultOptions: CookieAttributes = {
  expires: 365,
  secure: true,
  sameSite: "strict",
};

/**
 * @paramname The name of the cookie.
 * @returns The value of the cookie.
 */
export function getCookie(name: ICookieKey): string {
  const value = Cookies.get(name);

  if (!value) {
    throw new Error(`Failed to get value for the cookie "${name}".`);
  }

  return value;
}

/**
 * @param name The name of the cookie.
 * @param value The value of the cookie.
 * @param options Cookie options.
 */
export function setCookie(
  name: ICookieKey,
  value: string,
  options: Partial<CookieAttributes> = defaultOptions
) {
  const finalOptions = options
    ? { ...defaultOptions, ...options }
    : defaultOptions;
  const result = Cookies.set(name, value, finalOptions);

  if (!result) {
    throw new Error(`Failed to set the value for the cookie "${name}"`);
  }
}
