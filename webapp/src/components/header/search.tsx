import { useRouter } from "next/router";
import { useCallback, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSearch } from "@fortawesome/free-solid-svg-icons";
import { Button, Input } from "antd";

import { onParentBlur } from "#lib/dom";
export function Search() {
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
