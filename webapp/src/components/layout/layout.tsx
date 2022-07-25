import { type ReactNode } from "react";

import styles from "./layout.module.scss";
import { AccountNavigation } from "./account-nav";
import { GlobalSearch } from "./search";

import { LinkInternal } from "#components/links";
import { ListItem, ListUnordered } from "#components/lists";

export interface ILayoutProps {
  children: ReactNode;
}

export function Layout({ children }: ILayoutProps) {
  return (
    <>
      <header className={styles.header}>
        <nav className={styles.nav}>
          <ListUnordered className={styles.list}>
            <ListItem className={styles.logo}>
              {/* @TODO: site logo component */}
              <LinkInternal href="/">Horahora</LinkInternal>
            </ListItem>
            <ListItem className={styles.search}>
              <GlobalSearch />
            </ListItem>
            <AccountNavigation />
          </ListUnordered>
        </nav>
      </header>

      <main className={styles.main}>{children}</main>

      <footer className={styles.footer}>
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
