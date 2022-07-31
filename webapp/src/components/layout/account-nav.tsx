import {
  faArchive,
  faKey,
  faSignOutAlt,
  faUser,
  faUpload,
  faArrowRightToBracket,
  faUserPlus,
} from "@fortawesome/free-solid-svg-icons";
import { useState } from "react";
import clsx from "clsx";

import { ThemeSwitcher } from "./theme-switcher";

import { onParentBlur } from "#lib/dom";
import { ProfileURL } from "#lib/urls";
import { useAccount } from "#hooks";
import { LinkInternal } from "#components/links";
import { ListItem, ListUnordered } from "#components/lists";
import { Button } from "#components/buttons";

// eslint-disable-next-line
import styles from "./account-nav.module.scss";

export function AccountNavigation() {
  const [isOpen, switchOpen] = useState(false);
  const { account, isInProgress, isRegistered, isAdmin, logout } = useAccount();

  const className = clsx(
    styles.block,
    isOpen && styles.block_open,
    isInProgress && styles.block_loading
  );

  return (
    <ListItem
      className={className}
      onBlur={onParentBlur(() => {
        switchOpen(false);
      })}
    >
      <Button
        className={styles.switch}
        iconID={faUser}
        onClick={() => {
          switchOpen(!isOpen);
        }}
      >
        Account
      </Button>
      <ListUnordered className={styles.list}>
        {!isRegistered ? (
          <>
            <ListItem className={styles.item}>
              <LinkInternal
                className={styles.link}
                href="/authentication/login"
                iconID={faArrowRightToBracket}
              >
                Login
              </LinkInternal>
            </ListItem>

            <ListItem className={styles.item}>
              <LinkInternal
                className={styles.link}
                href="/authentication/register"
                iconID={faUserPlus}
              >
                Register
              </LinkInternal>
            </ListItem>
          </>
        ) : (
          <>
            <ListItem className={styles.item}>
              <LinkInternal
                className={styles.link}
                iconID={faUser}
                // @ts-expect-error @TODO: better hook
                href={new ProfileURL(account.userID)}
              >
                Profile page
              </LinkInternal>
            </ListItem>

            <ListItem className={styles.item}>
              <LinkInternal
                className={styles.link}
                iconID={faUpload}
                href="/account/upload"
              >
                Upload
              </LinkInternal>
            </ListItem>

            {isAdmin && (
              <ListItem className={styles.item}>
                <LinkInternal
                  className={styles.link}
                  iconID={faArchive}
                  href="/account/archives"
                >
                  Archives
                </LinkInternal>
              </ListItem>
            )}

            <ListItem className={styles.item}>
              <ThemeSwitcher />
            </ListItem>

            <ListItem className={styles.item}>
              <LinkInternal
                className={styles.link}
                iconID={faKey}
                href="/authentication/password-reset"
              >
                Reset Password
              </LinkInternal>
            </ListItem>

            {isAdmin && (
              <ListItem className={styles.item}>
                <LinkInternal
                  className={styles.link}
                  iconID={faArchive}
                  href="/account/administrator/audits"
                >
                  Audit Logs
                </LinkInternal>
              </ListItem>
            )}

            <ListItem className={styles.item}>
              <Button
                iconID={faSignOutAlt}
                onClick={async () => {
                  await logout();
                }}
              >
                Logout
              </Button>
            </ListItem>
          </>
        )}
      </ListUnordered>
    </ListItem>
  );
}
