import {useCallback, useState} from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArchive,
  faBars,
  faSearch,
  faSignOutAlt,
  faUser,
} from "@fortawesome/free-solid-svg-icons";
import { useHistory } from 'react-router-dom';
import { Button, Dropdown, Input, Menu } from "antd";

function Search() {
    const [redirectVal, setRedirectVal] = useState(null);

    let onSubmit = useCallback((e) => {
        e.preventDefault();
        const category = document.getElementById('category').value;
        const order = document.querySelector('input[name="order"]:checked').value;
        const search = document.querySelector('input[name="search"]').value;

        setRedirectVal(`/?search=${search}&order=${order}&category=${category}`)
    }, []);

    const history = useHistory();
    if (redirectVal){
        history.push(redirectVal);
        setRedirectVal(null);
    }

  return (
    <>
      <form onSubmit={onSubmit} className="w-full max-w-sm z-50" onMouseEnter={showModal} onMouseLeave={hideModal}>
        <Input
          name="search"
          size="large"
          placeholder="Search"
          prefix={
            <FontAwesomeIcon className="mr-1 text-gray-400" icon={faSearch} />
          }
        />
              <div id="hidden-search-modal" className="absolute bg-white w-full max-w-sm p-5 space-y-3 invisible">
                  <h1>SEARCH OPTIONS</h1>
                  Order by
                  <select name="category" id="category">
                      <option value="upload_date">upload date</option>
                      <option value="rating">rating</option>
                      <option value="views">views</option>
                  </select>
                  <br></br>
                    <input type="radio" id="desc" name="order" value="desc"></input>
                          <label htmlFor="desc">Desc</label>
                    <input type="radio" id="asc" name="order" value="asc"></input>
                            <label htmlFor="asc">Asc</label>
                  <br></br>
                  <Button block type="primary" htmlType="submit" size="large">Submit</Button>
              </div>
      </form>
    </>
  );
}

function LoggedInUserNav(props) {
  const { userData } = props;

  let menu = (
    <Menu>
      <Menu.Item key="profile" icon={<FontAwesomeIcon icon={faUser} />}>
        <Link to={`/users/${userData.userID}`}>Profile page</Link>
      </Menu.Item>
      <Menu.Item
        key="archive-requests"
        icon={<FontAwesomeIcon icon={faArchive} />}
      >
        <Link to="/archive-requests">Archive Requests</Link>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item
        key="logout"
        icon={<FontAwesomeIcon className="text-red-600" icon={faSignOutAlt} />}
      >
        <Link to="/logout">Logout</Link>
      </Menu.Item>
    </Menu>
  );

  return (
    <>
      <Dropdown overlay={menu} placement="bottomRight" trigger={["click"]}>
        <Button>
          {userData.username}
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

function UserNav(props) {
  const { userData } = props;

  if (userData && userData.username) {
    return <LoggedInUserNav userData={userData} />;
  } else {
    return <LoggedOutUserNav />;
  }
}

function showModal(){
    console.log(document.getElementById('hidden-search-modal'));
    document.getElementById('hidden-search-modal').style.visibility = 'visible';
}

function hideModal(){
    document.getElementById('hidden-search-modal').style.visibility = 'hidden';
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
