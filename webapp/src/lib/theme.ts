import { getCookie, setCookie } from "../store/cookie";

export const THEME = {
  DARK: "dark",
  LIGHT: "light",
} as const;

export type ITheme = typeof THEME[keyof typeof THEME];

export const DEFAULT_THEME = THEME.DARK;

export function getTheme(): ITheme {
  let currentTheme: string | undefined = undefined;

  try {
    currentTheme = getCookie("theme");
  } catch (error) {
    currentTheme = document.documentElement.dataset.theme;
    setCookie("theme", currentTheme as ITheme);
  }

  return currentTheme as ITheme;
}

export function setTheme(theme: ITheme): ITheme {
  document.documentElement.dataset.theme = theme;
  setCookie("theme", theme);

  return theme;
}
