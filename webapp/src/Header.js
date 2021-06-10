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

function UserNav() {
  return <>
    todo: user nav
  </>
}

function Header() {
  return <nav className='h-16 bg-white shadow flex justify-center'>
    <div className="max-w-screen-lg w-screen flex justify-around items-center">
      <div className="flex justify-start">
        <Link to='/'>Horahora</Link>
      </div>
      <div className="flex-grow flex justify-center">
        <Search />
      </div>
      <div className="">
        <UserNav />
      </div>
    </div>
  </nav>
}

export default Header;
