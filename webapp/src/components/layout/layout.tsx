import { type ReactNode } from "react";

import styles from "./layout.module.scss";
import { AccountNavigation } from "./account-nav";

import { LinkInternal } from "#components/links";
import { Search } from "#components/inputs";
import { ListItem, ListUnordered } from "#components/lists";

export interface ILayoutProps {
  children: ReactNode;
  userData?: Record<string, unknown>;
  dataless?: boolean;
}

export function Layout({ userData, dataless, children }: ILayoutProps) {
  return (
    <>
      <header className={styles.header}>
        <nav className={styles.nav}>
          <ListUnordered >
            <ListItem>
              {/* @TODO: site logo component */}
              <LinkInternal className={styles.logo} href="/">
              Horahora
            </LinkInternal></ListItem>
          </ListUnordered>
          <div className="flex justify-start flex-grow-0">


          </div>
          <Search />
          {!dataless && (
            <div className="flex-grow-0 ml-auto">
              <AccountNavigation userData={userData} />
            </div>
          )}
        </nav>
      </header>

      <main className={styles.main}>{children}</main>

      <footer className={styles.block}>
        <LinkInternal className={styles.link} href="/privacy-policy">
          Privacy Policy
        </LinkInternal>
        <LinkInternal className={styles.link} href="/terms-of-service">
          Terms of Service
        </LinkInternal>
      </footer>
    </>
  );
}
