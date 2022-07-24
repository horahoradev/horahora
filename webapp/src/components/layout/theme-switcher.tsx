import { useEffect, useState } from "react";
import { faSun } from "@fortawesome/free-solid-svg-icons";

import {
  DEFAULT_THEME,
  getTheme,
  setTheme,
  THEME,
  type ITheme,
} from "#lib/theme";
import { Button } from "#components/buttons";

/**
 * @TODO more standalone markup.
 */
export function ThemeSwitcher() {
  const [currentTheme, switchTheme] = useState<ITheme>(DEFAULT_THEME);

  // initialize component on client
  useEffect(() => {
    switchTheme(getTheme());
  }, []);

  return (
    <Button
      iconID={faSun}
      onClick={() => {
        const nextTheme =
          currentTheme === THEME.DARK ? THEME.LIGHT : THEME.DARK;
        switchTheme(setTheme(nextTheme));
      }}
    ></Button>
  );
}
