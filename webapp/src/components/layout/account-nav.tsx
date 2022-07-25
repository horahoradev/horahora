import {
  faArchive,
  faKey,
  faSignOutAlt,
  faUser,
  faUpload,
  faArrowRightToBracket,
  faUserPlus,
} from "@fortawesome/free-solid-svg-icons";
import { useState, useEffect } from "react";
import clsx from "clsx";

import { ThemeSwitcher } from "./theme-switcher";

import { getUserdata } from "#api/index";
import { LinkInternal } from "#components/links";
import { UserRank } from "#api/types";
import { ListItem, ListUnordered } from "#components/lists";
import { Button } from "#components/buttons";
import { type ILoggedInUserData } from "#codegen/schema/001_interfaces";
import { onParentBlur } from "#lib/dom";
import { ProfileURL } from "#lib/urls";

// eslint-disable-next-line
import styles from "./account-nav.module.scss";

export function AccountNavigation() {
  const [isLoading, switchLoading] = useState(true);
  const [isOpen, switchOpen] = useState(false);
  const [userData, setUserData] = useState<ILoggedInUserData>();

  const className = clsx(
    styles.block,
    isOpen && styles.block_open,
    isLoading && styles.block_loading
  );
  const isRegistered = Boolean(userData && userData.username);
  const isAdmin = userData?.rank === UserRank.ADMIN;

  useEffect(() => {
    (async () => {
      try {
        let userData = await getUserdata();
        setUserData(userData);
      } finally {
        switchLoading(false);
      }
    })();
  }, []);

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
                // @ts-expect-error figure `userData` shape
                href={new ProfileURL(userData.userID)}
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
                  href="/account/archive-requests"
                >
                  Archive Requests
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
              <LinkInternal
                className={styles.link}
                iconID={faSignOutAlt}
                href="/authentication/logout"
              >
                Logout
              </LinkInternal>
            </ListItem>
          </>
        )}
      </ListUnordered>
    </ListItem>
  );
}
