import { useRouter } from "next/router";
import { useCallback, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArchive,
  faBars,
  faSearch,
  faKey,
  faSignOutAlt,
  faUser,
  faUpload,
} from "@fortawesome/free-solid-svg-icons";
import { Button, Dropdown, Input, Menu } from "antd";

import { ThemeSwitcher } from "./theme-switcher";
import { LinkInternal } from "./links/internal";

import { UserRank } from "#api/types";
import { onParentBlur } from "#lib/dom";
import { Icon } from "#components/icons";

interface IHeaderProps extends Record<string, unknown> {}
export function Header({ userData, dataless }: IHeaderProps) {
  return (
    <nav className="h-16 bg-white dark:bg-gray-800 shadow flex justify-center">
      <div className="max-w-screen-lg w-screen flex justify-start items-center gap-x-4 mx-4">
        <div className="flex justify-start flex-grow-0">
          {/* @TODO: site logo component */}
          <LinkInternal className="text-2xl font-black text-blue-500" href="/">
            Horahora
          </LinkInternal>
        </div>
        <Search />
        {!dataless && (
          <div className="flex-grow-0 ml-auto">
            <UserNav userData={userData} />
          </div>
        )}
      </div>
    </nav>
  );
}

function Search() {
  const router = useRouter();
  const [redirectVal, setRedirectVal] = useState<string | null>(null);
  const [isFocused, switchFocus] = useState(false);

  let handleSearch = useCallback((e) => {
    e.preventDefault();
    const category = (document.getElementById("category") as HTMLSelectElement)
      .value;
    const order = (
      document.querySelector('input[name="order"]:checked') as HTMLInputElement
    ).value;
    const search = (
      document.querySelector('input[name="search"]') as HTMLInputElement
    ).value;

    const params = new URLSearchParams([
      ["category", category],
      ["order", order],
      ["search", search],
    ]);

    setRedirectVal(`/?${params.toString()}`);
  }, []);

  if (redirectVal) {
    router.push(redirectVal);
    setRedirectVal(null);
  }

  return (
    <form
      onSubmit={handleSearch}
      className="flex-grow flex flex-col w-full max-w-sm"
      onBlur={onParentBlur(() => {
        switchFocus(false);
      })}
    >
      <Input
        name="search"
        size="large"
        placeholder="Search"
        prefix={
          <FontAwesomeIcon
            className="mr-1 text-gray-400 max-h-4"
            icon={faSearch}
          />
        }
        onFocus={() => {
          switchFocus(true);
        }}
        onBlur={(event) => {
          event.preventDefault();
        }}
      />
      <div
        className={
          isFocused
            ? "absolute top-14 z-10 text-black bg-white dark:text-white dark:bg-black  w-full max-w-sm p-4 visible opacity-1 duration-250 transition-opacity transition-visibility"
            : "absolute top-14 z-10 text-black bg-white dark:text-white dark:bg-black w-full max-w-sm p-4 invisible opacity-0 duration-250 transition-opacity transition-visibility"
        }
        tabIndex={0}
      >
        <label
          htmlFor="category"
          className="text-black dark:text-white text-base"
        >
          Order by:{" "}
        </label>
        <select
          name="category"
          className="bg-white dark:bg-black"
          id="category"
        >
          <option value="upload_date">upload date</option>
          <option value="rating">rating</option>
          <option value="views">views</option>
          <option value="my_ratings">my ratings</option>
        </select>
        <br />
        <input
          type="radio"
          id="desc"
          name="order"
          defaultChecked={true}
          value="desc"
        ></input>
        <label htmlFor="desc">Desc</label>
        <input type="radio" id="asc" name="order" value="asc"></input>
        <label htmlFor="asc">Asc</label>
        <br />
        <Button block type="primary" htmlType="submit" size="large">
          Submit
        </Button>
      </div>
    </form>
  );
}

interface IUserNav extends Record<string, unknown> {}

function UserNav(props: IUserNav) {
  const { userData } = props;

  // @ts-expect-error figure `userData` shape
  if (userData && userData.username && userData.rank === UserRank.ADMIN) {
    return <LoggedInAdminNav userData={userData} />;
    // @ts-expect-error figure `userData` shape
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

      <Menu.Item
        key="password-reset"
      >
        <LinkInternal iconID={faKey} href="/password-reset">Reset Password</LinkInternal>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item
        key="logout"
      >
        <LinkInternal iconID={faSignOutAlt} href="/logout">Logout</LinkInternal>
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
      <Menu.Item
        key="profile"
      >
        <LinkInternal
          iconID={faUser}
          // @ts-expect-error figure `userData` shape
          href={`/users/${userData.userID}`}
        >
          Profile page
        </LinkInternal>
      </Menu.Item>
      <Menu.Divider />

      <Menu.Item
        key="upload"
        className="flex flex-row gap-2"
      >
        <LinkInternal iconID={faUpload} href="/account/upload">Upload</LinkInternal>
      </Menu.Item>
      <Menu.Divider />

      <Menu.Item
        key="archive-requests"
      >
        <LinkInternal iconID={faArchive} href="/archive-requests">Archive Requests</LinkInternal>
      </Menu.Item>
      <Menu.Divider />

      <ThemeSwitcher />

      <Menu.Divider />

      <Menu.Item
        key="password-reset"
      >
        <LinkInternal iconID={faKey} href="/password-reset">Password Reset</LinkInternal>
      </Menu.Item>

      <Menu.Divider />
      <Menu.Item
        key="audits"
      >
        <LinkInternal iconID={faArchive} href="/audits">Audit Logs</LinkInternal>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item
        key="logout"
      >
        <LinkInternal iconID={faSignOutAlt} href="/logout">Logout</LinkInternal>
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
      <LinkInternal href="/login">Login</LinkInternal>
      <LinkInternal href="/register">Register</LinkInternal>
    </>
  );
}
