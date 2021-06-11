import { useCallback } from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSearch } from "@fortawesome/free-solid-svg-icons";
import { Input } from "antd";

function Search() {
  let onSubmit = useCallback((e) => {
    e.preventDefault();
    // TODO(ivan): Search
  }, []);

  return (
    <>
      <form onSubmit={onSubmit} className="w-full max-w-sm">
        <Input
          name="username"
          size="large"
          placeholder="Search"
          prefix={
            <FontAwesomeIcon className="mr-1 text-gray-400" icon={faSearch} />
          }
        />
      </form>
    </>
  );
}

function LoggedInUserNav(props) {
  const { userData } = props;

  return (
    <>
      <Link className="mr-4" to="/archive-requests">
        Archive Requests
      </Link>
      <Link to={`/users/${userData.UserID}`}>{userData.Username}</Link>
    </>
  );
}

function LoggedOutUserNav() {
  return (
    <>
      <Link to="/login">Login</Link>
    </>
  );
}

function UserNav(props) {
  const { userData } = props;

  if (userData.Username) {
    return <LoggedInUserNav userData={userData} />;
  } else {
    return <LoggedOutUserNav />;
  }
}

function Header(props) {
  const { userData, dataless } = props;

  return (
    <nav className="h-16 bg-white shadow flex justify-center">
      <div className="max-w-screen-lg w-screen flex justify-start items-center mx-4">
        <div className="flex justify-start flex-grow-0">
          <Link to="/" className="text-xl">
            Horahora
          </Link>
        </div>
        <div className="flex-grow flex mx-8">
          <Search />
        </div>
        {!dataless && (
          <div className="flex-grow-0">
            <UserNav userData={userData} />
          </div>
        )}
      </div>
    </nav>
  );
}

export default Header;
