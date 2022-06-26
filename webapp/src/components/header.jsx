import { useCallback, useState } from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArchive,
  faBars,
  faSearch,
  faKey,
  faSignOutAlt,
  faUser,
  faSun,
} from "@fortawesome/free-solid-svg-icons";
import { useHistory } from "react-router-dom";
import { Switch, Button, Dropdown, Input, Menu } from "antd";

import { UserRank } from "../api/types";
import { onParentBlur } from "../lib/dom";

export function Header({ userData, dataless }) {
  return (
    <nav className="h-16 bg-white dark:bg-gray-800 shadow flex justify-center">
      <div className="max-w-screen-lg w-screen flex justify-start items-center gap-x-4 mx-4">
        <div className="flex justify-start flex-grow-0">
          <Link to="/" className="text-2xl font-black text-blue-500">
            Horahora
          </Link>
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
  const [redirectVal, setRedirectVal] = useState(null);
  const [isFocused, switchFocus] = useState(false);

  let handleSearch = useCallback((e) => {
    e.preventDefault();
    const category = document.getElementById("category").value;
    const order = document.querySelector('input[name="order"]:checked').value;
    const search = document.querySelector('input[name="search"]').value;

    const params = new URLSearchParams([
      ["category", category],
      ["order", order],
      ["search", search],
    ]);

    setRedirectVal(`/?${params.toString()}`);
  }, []);

  const history = useHistory();
  if (redirectVal) {
    history.push(redirectVal);
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
          <FontAwesomeIcon className="mr-1 text-gray-400" icon={faSearch} />
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

function UserNav(props) {
  const { userData } = props;

  if (userData && userData.username && userData.rank === UserRank.ADMIN) {
    return <LoggedInAdminNav userData={userData} />;
  } else if (userData && userData.username) {
    return <LoggedInUserNav userData={userData} />;
  } else {
    return <LoggedOutUserNav />;
  }
}

function LoggedInUserNav(props) {
  const { userData } = props;

  let menu = (
    <Menu className="bg-white dark:bg-black">
      <Menu.Item
        key="profile"
        icon={
          <FontAwesomeIcon
            className="text-black dark:text-white"
            icon={faUser}
          />
        }
      >
        <Link
          className="text-black dark:text-white dark:hover:text-black"
          to={`/users/${userData.userID}`}
        >
          Profile page
        </Link>
      </Menu.Item>
      <Menu.Divider />

      <Menu.Item
        key="password-reset"
        icon={
          <FontAwesomeIcon
            className="text-black dark:text-white"
            icon={faKey}
          />
        }
      >
        <Link className="text-black dark:text-white" to="/password-reset">
          Reset Password
        </Link>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item
        key="logout"
        icon={<FontAwesomeIcon className="text-red-600" icon={faSignOutAlt} />}
      >
        <Link className="text-black dark:text-white" to="/logout">
          Logout
        </Link>
      </Menu.Item>
    </Menu>
  );

  return (
    <>
      <Dropdown overlay={menu} placement="bottomRight" trigger={["click"]}>
        <Button>
          <b className="text-blue-500">{userData.username}</b>
          <FontAwesomeIcon className="text-xs ml-2" icon={faBars} />
        </Button>
      </Dropdown>
    </>
  );
}

function LoggedInAdminNav(props) {
  const { userData } = props;

  let toggleDarkmode = function () {};

  let menu = (
    <Menu className="bg-white dark:bg-black">
      <Menu.Item
        key="profile"
        icon={
          <FontAwesomeIcon
            className="text-black dark:text-white"
            icon={faUser}
          />
        }
      >
        <Link
          className="text-black dark:text-white dark:hover:text-black"
          to={`/users/${userData.userID}`}
        >
          Profile page
        </Link>
      </Menu.Item>

      <Menu.Divider />
      <Menu.Item
        key="archive-requests"
        icon={
          <FontAwesomeIcon
            className="text-black dark:text-white"
            icon={faArchive}
          />
        }
      >
        <Link className="text-black dark:text-white dark:hover:text-black" to="/archive-requests">
          Archive Requests
        </Link>
      </Menu.Item>
      <Menu.Divider />

      <Menu.Item
        key="darkmode"
        icon={
          <FontAwesomeIcon
            className="text-black dark:text-white"
            icon={faSun}
          />
        }
      >
        <Switch
          className="text-black dark:text-white "
          onChange={() => toggleDarkmode()}
        ></Switch>
      </Menu.Item>

      <Menu.Divider />

      <Menu.Item
        key="password-reset"
        icon={
          <FontAwesomeIcon
            className="text-black dark:text-white"
            icon={faKey}
          />
        }
      >
        <Link className="text-black dark:text-white dark:hover:text-black" to="/password-reset">
          Password Reset
        </Link>
      </Menu.Item>

      <Menu.Divider />
      <Menu.Item
        key="audits"
        icon={
          <FontAwesomeIcon
            className="text-black dark:text-white"
            icon={faArchive}
          />
        }
      >
        <Link className="text-black dark:text-white dark:hover:text-black" to="/audits">
          Audit Logs
        </Link>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item
        key="logout"
        icon={<FontAwesomeIcon className="text-red-600" icon={faSignOutAlt} />}
      >
        <Link className="text-black dark:text-white" to="/logout">
          Logout
        </Link>
      </Menu.Item>
    </Menu>
  );

  return (
    <>
      <Dropdown overlay={menu} placement="bottomRight" trigger={["click"]}>
        <Button>
          <b className="text-blue-500">{userData.username}</b>
          <FontAwesomeIcon className="text-xs ml-2" icon={faBars} />
        </Button>
      </Dropdown>
    </>
  );
}

function LoggedOutUserNav() {
  return (
    <>
      <Link to="/login">
        <Button>Login</Button>
      </Link>
      <Link className="ml-2" to="/register">
        <Button type="primary">Register</Button>
      </Link>
    </>
  );
}
