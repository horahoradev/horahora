import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArchive,
  faBars,
  faKey,
  faSignOutAlt,
  faUser,
  faUpload,
} from "@fortawesome/free-solid-svg-icons";
import { Dropdown, Menu } from "antd";
import { useState, useEffect } from "react";

import styles from "./account-nav.module.scss";

import { getUserdata } from "#api/index";
import { ThemeSwitcher } from "#components/theme-switcher";
import { LinkInternal } from "#components/links";
import { UserRank } from "#api/types";
import { ListItem, ListUnordered } from "#components/lists";
import { Button } from "#components/buttons";
import { type ILoggedInUserData } from "#codegen/schema/001_interfaces";

export function AccountNavigation() {
  const [userData, setUserData] = useState<ILoggedInUserData>();
  const isRegistered = Boolean(userData && userData.username);
  const isAdmin = userData?.rank === UserRank.ADMIN;

  useEffect(() => {
    let ignore = false;

    (async () => {
      let userData = await getUserdata();
      if (!ignore) setUserData(userData);
    })();

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <ListUnordered className={styles.block}>
      <ListItem>
        <b className="text-blue-500">{userData?.username}</b>
        <FontAwesomeIcon className="max-h-4 text-xs ml-2" icon={faBars} />
        <Button>Account</Button>
        <ListUnordered>
          {!isRegistered ? (
            <>
              <ListItem>
                <LinkInternal href="/authentication/login">Login</LinkInternal>
              </ListItem>
              <ListItem>
                <LinkInternal href="/authentication/register">
                  Register
                </LinkInternal>
              </ListItem>
            </>
          ) : (
            <>
              <ListItem>
                <LinkInternal iconID={faUpload} href="/account/upload">
                  Upload
                </LinkInternal>
              </ListItem>
              <ListItem>
                <LinkInternal
                  iconID={faUser}
                  // @ts-expect-error figure `userData` shape
                  href={`/users/${userData.userID}`}
                >
                  Profile page
                </LinkInternal>
              </ListItem>
              <ListItem>
                <LinkInternal
                  iconID={faKey}
                  href="/authentication/password-reset"
                >
                  Reset Password
                </LinkInternal>
              </ListItem>
              <ListItem>
                <LinkInternal
                  iconID={faSignOutAlt}
                  href="/authentication/logout"
                >
                  Logout
                </LinkInternal>
              </ListItem>
            </>
          )}
        </ListUnordered>
      </ListItem>
    </ListUnordered>
  );
}

interface ILoggedInAdminNav extends Record<string, unknown> {}

function LoggedInAdminNav(props: ILoggedInAdminNav) {
  const { userData } = props;

  let menu = (
    <Menu className="bg-white dark:bg-black">
      <Menu.Item key="profile">
        <LinkInternal
          iconID={faUser}
          // @ts-expect-error figure `userData` shape
          href={`/users/${userData.userID}`}
        >
          Profile page
        </LinkInternal>
      </Menu.Item>
      <Menu.Divider />

      <Menu.Item key="upload" className="flex flex-row gap-2">
        <LinkInternal iconID={faUpload} href="/account/upload">
          Upload
        </LinkInternal>
      </Menu.Item>
      <Menu.Divider />

      <Menu.Item key="archive-requests">
        <LinkInternal iconID={faArchive} href="/account/archive-requests">
          Archive Requests
        </LinkInternal>
      </Menu.Item>
      <Menu.Divider />

      <ThemeSwitcher />

      <Menu.Divider />

      <Menu.Item key="password-reset">
        <LinkInternal iconID={faKey} href="/authentication/password-reset">
          Password Reset
        </LinkInternal>
      </Menu.Item>

      <Menu.Divider />
      <Menu.Item key="audits">
        <LinkInternal iconID={faArchive} href="/account/administrator/audits">
          Audit Logs
        </LinkInternal>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item key="logout">
        <LinkInternal iconID={faSignOutAlt} href="/authentication/logout">
          Logout
        </LinkInternal>
      </Menu.Item>
    </Menu>
  );

  return (
    <>
      <Dropdown overlay={menu} placement="bottomRight" trigger={["click"]}>
        <Button className="flex flex-row flex-nowrap items-center">
          {/* @ts-expect-error figure `userData` shape */}
          <b className="text-blue-500">{userData.username}</b>
          <FontAwesomeIcon
            className="text-xs max-h-4 ml-2 text-black dark:text-white dark:hover:text-black"
            icon={faBars}
          />
        </Button>
      </Dropdown>
    </>
  );
}
