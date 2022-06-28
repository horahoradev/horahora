import { useEffect, useState } from "react";
import { Switch, Menu } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSun } from "@fortawesome/free-solid-svg-icons";

import { DEFAULT_THEME, getTheme, setTheme } from "../lib/theme";

/**
 * @TODO more standalone markup.
 */
export function ThemeSwitcher() {
  const [currentTheme, switchTheme] = useState(DEFAULT_THEME);

  // initialize component on client
  useEffect(() => {
    switchTheme(getTheme);
  }, []);

  return (
    <Menu.Item
      key="darkmode"
      className="bg-white dark:bg-black"
      icon={
        <FontAwesomeIcon className="max-h-4 text-black dark:text-white" icon={faSun} />
      }
    >
      <Switch
        className="text-black dark:text-white"
        onChange={() => {
          const nextTheme = currentTheme === "dark" ? "light" : "dark";
          switchTheme(setTheme(nextTheme));
        }}
      ></Switch>
    </Menu.Item>
  );
}
