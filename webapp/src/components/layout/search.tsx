import { useRouter } from "next/router";
import { useCallback, useState } from "react";
import {
  faSearch,
} from "@fortawesome/free-solid-svg-icons";
import { Button, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import { onParentBlur } from "#lib/dom";

import clsx from "clsx";
import styles from "./search.module.scss";

export function GlobalSearch() {
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

  const className = clsx(
    isFocused && styles.form__focus,
    !isFocused && styles.form__blur
  );

  return (
    <form
      onSubmit={handleSearch}
      className={styles.form}
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
            className={styles.form__search}
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
        className={className}
        tabIndex={0}
      >
        <label
          htmlFor="category"
          className={styles.form__label}
        >
          Order by:{" "}
        </label>
        <select
          name="category"
          className={styles.form__select}
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
