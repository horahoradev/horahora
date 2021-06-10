import { useRef, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearch } from '@fortawesome/free-solid-svg-icons';

function Search() {
  let searchRef = useRef();

  let searchInputFocus = useCallback(() => {
    searchRef.current?.focus();
  }, [searchRef]);

  let onSubmit = useCallback((e) => {
    e.preventDefault();
    // TODO(ivan): Search
  }, []);

  // TODO(ivan): Expanding search input (with max-width)
  return <>
    <form onSubmit={onSubmit}>
      <div className="border rounded bg-white focus-within:ring" onMouseUp={searchInputFocus}>
        <FontAwesomeIcon htmlFor="input" className="mx-2 text-gray-400" icon={faSearch} />
        <input
          ref={searchRef}
          name="input"
          className="py-1 pr-3 focus:outline-none"
          type='text'
          placeholder='Search'
        />
      </div>
    </form>
  </>
}

function LoggedInUserNav(props) {
  const { userData } = props;

  return <>
    <Link className="mr-4" to="/archive-requests">Archive Requests</Link>
    <Link to={`/users/${userData.UserID}`}>{userData.Username}</Link>
  </>
}

function LoggedOutUserNav() {
  return <>
    <Link to="/login">Login</Link>
  </>
}

function UserNav(props) {
  const { userData } = props;

  if (userData.Username) {
    return <LoggedInUserNav userData={userData} />
  } else {
    return <LoggedOutUserNav />
  }
}

function Header(props) {
  const { userData } = props;

  return <nav className='h-16 bg-white shadow flex justify-center'>
    <div className="max-w-screen-lg w-screen flex justify-around items-center mx-4">
      <div className="flex justify-start">
        <Link to='/'>Horahora</Link>
      </div>
      <div className="flex-grow flex justify-center">
        <Search />
      </div>
      <div className="">
        <UserNav userData={userData} />
      </div>
    </div>
  </nav>
}

export default Header;
