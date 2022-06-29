import { getCookie, setCookie } from "../store/cookie";

/**
 * @typedef {"dark" | "light"} ITheme
 */

export const DEFAULT_THEME = "dark";

/**
 * @returns {ITheme}
 */
export function getTheme() {
  /**
   * @typedef {ITheme}
   */
  let currentTheme = undefined;

  try {
    currentTheme = getCookie("theme");
  } catch (error) {
    currentTheme = document.documentElement.dataset.theme;
    setCookie("theme", currentTheme);
  }

  return currentTheme;
}

/**
 * @param {ITheme} theme
 */
export function setTheme(theme) {
  document.documentElement.dataset.theme = theme;
  setCookie("theme", theme);
  
  return theme;
}
