import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArchive,
  faBars,
  faKey,
  faSignOutAlt,
  faUser,
  faUpload,
} from "@fortawesome/free-solid-svg-icons";
import { Button, Dropdown, Menu } from "antd";
import { useState, useEffect } from "react";
import { useRouter } from "next/router";

import { getUserdata } from "#api/index";
import { ThemeSwitcher } from "#components/theme-switcher";
import { LinkInternal } from "#components/links";
import { UserRank } from "#api/types";

export function AccountNavigation() {
  const router = useRouter();
  const [userData, setUserData] = useState<Record<string, unknown>>();

  useEffect(() => {
    let ignore = false;

    if (router.pathname.startsWith("/authentication")) {
      return () => {
        ignore = true;
      };
    }

    (async () => {
      let userData = await getUserdata();
      if (!ignore) setUserData(userData);
    })();

    return () => {
      ignore = true;
    };
  }, []);

  const isRegistered = Boolean(userData && userData.username);
  const isAdmin = userData?.rank === UserRank.ADMIN;

  if (isRegistered && isAdmin) {
    return <LoggedInAdminNav userData={userData} />;
  } else if (userData && userData.username) {
    return <LoggedInUserNav userData={userData} />;
  } else {
    return <LoggedOutUserNav />;
  }
}

interface ILoggedInUserNav extends Record<string, unknown> {}

function LoggedInUserNav(props: ILoggedInUserNav) {
  const { userData } = props;

  let menu = (
    <Menu className="bg-white dark:bg-black">
      <Menu.Item key="upload">
        <LinkInternal iconID={faUpload} href="/account/upload">
          Upload
        </LinkInternal>
      </Menu.Item>
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

      <Menu.Item key="password-reset">
        <LinkInternal iconID={faKey} href="/authentication/password-reset">
          Reset Password
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
        <Button>
          {/* @ts-expect-error figure `userData` shape */}
          <b className="text-blue-500">{userData.username}</b>
          <FontAwesomeIcon className="max-h-4 text-xs ml-2" icon={faBars} />
        </Button>
      </Dropdown>
    </>
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

function LoggedOutUserNav() {
  return (
    <>
      <LinkInternal href="/authentication/login">Login</LinkInternal>
      <LinkInternal href="/authentication/register">Register</LinkInternal>
    </>
  );
}
