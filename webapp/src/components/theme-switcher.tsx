import { useEffect, useState } from "react";
import { Switch, Menu } from "antd";
import { faSun } from "@fortawesome/free-solid-svg-icons";

import {
  DEFAULT_THEME,
  getTheme,
  setTheme,
  THEME,
  type ITheme,
} from "#lib/theme";
import { Icon } from "#components/icons";

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
    <Menu.Item
      key="darkmode"
      className="bg-white dark:bg-black"
      icon={<Icon icon={faSun} />}
    >
      <Switch
        className="text-black dark:text-white"
        onChange={() => {
          const nextTheme =
            currentTheme === THEME.DARK ? THEME.LIGHT : THEME.DARK;
          switchTheme(setTheme(nextTheme));
        }}
      ></Switch>
    </Menu.Item>
  );
}
