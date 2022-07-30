import { type ReactNode } from "react";

import { Heading } from "#components/headings";

interface IPageProps {
  title?: ReactNode;
  children: ReactNode;
}

/**
 * If no `title` prop is passed, then it's up to the page
 * to construct and style the page markup.
 */
export function Page({ title, children }: IPageProps) {
  return (
    <>
      {!title ? (
        children
      ) : (
        <>
          <Heading level={1}>{title}</Heading>
          <section>{children}</section>
        </>
      )}
    </>
  );
}
