import { getCookie, setCookie } from "../store/cookie";

type ITheme = "dark" | "light";

export const DEFAULT_THEME = "dark";

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
