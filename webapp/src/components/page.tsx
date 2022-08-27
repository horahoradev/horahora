import { type ReactNode } from "react";
import Head from "next/head";

import { Heading } from "#components/headings";

// eslint-disable-next-line
import styles from "./page.module.scss";

interface IPageProps {
  title?: string;
  heading?: ReactNode;
  children: ReactNode;
}

/**
 * If no `title` prop is passed, then it's up to the page
 * to construct and style the page markup.
 */
export function Page({ title, heading = title, children }: IPageProps) {
  return (
    <>
      <Head>
        <title key={title}>{title}</title>
      </Head>
      {!heading ? (
        children
      ) : (
        <>
          <Heading level={1}>{heading}</Heading>
          <section className={styles.content}>{children}</section>
        </>
      )}
    </>
  );
}
