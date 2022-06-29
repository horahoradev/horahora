import Cookies from "js-cookie";

/**
 * @typedef {import("js-cookie").CookieAttributes} CookieAttributes
 */

/**
 * Only includes client-accessible cookies.
 * @typedef {"theme"} ICookieKey
 */

/**
 * @type {CookieAttributes}
 */
const defaultOptions = {
  expires: 365,
  secure: true,
  sameSite: "strict",
};

/**
 * @param {ICookieKey} name The name of the cookie.
 * @returns {string} The value of the cookie.
 */
export function getCookie(name) {
  const value = Cookies.get(name);

  if (!value) {
    throw new Error(`Failed to get value for the cookie "${name}".`);
  }

  return value;
}

/**
 * @param {ICookieKey} name The name of the cookie.
 * @param {string} value The value of the cookie.
 * @param {Partial<CookieAttributes>} options Cookie options.
 */
export function setCookie(name, value, options = defaultOptions) {
  const finalOptions = options
    ? { ...defaultOptions, ...options }
    : defaultOptions;
  const result = Cookies.set(name, value, finalOptions);

  if (!result) {
    throw new Error(`Failed to set the value for the cookie "${name}"`);
  }
}
