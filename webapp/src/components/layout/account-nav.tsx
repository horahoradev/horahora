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
            <ListItem>
              <LinkInternal
                href="/authentication/login"
                iconID={faArrowRightToBracket}
              >
                Login
              </LinkInternal>
            </ListItem>
            <ListItem>
              <LinkInternal href="/authentication/register" iconID={faUserPlus}>
                Register
              </LinkInternal>
            </ListItem>
          </>
        ) : (
          <>
            <ListItem>
              <LinkInternal iconID={faUpload} href="/account/upload">
                Upload
              </LinkInternal>
            </ListItem>
            <ListItem>
              <LinkInternal
                iconID={faUser}
                // @ts-expect-error figure `userData` shape
                href={`/users/${userData.userID}`}
              >
                Profile page
              </LinkInternal>
            </ListItem>
            {isAdmin && (
              <ListItem>
                <LinkInternal
                  iconID={faArchive}
                  href="/account/archive-requests"
                >
                  Archive Requests
                </LinkInternal>
              </ListItem>
            )}
            <ListItem>
              <ThemeSwitcher />
            </ListItem>
            <ListItem>
              <LinkInternal
                iconID={faKey}
                href="/authentication/password-reset"
              >
                Reset Password
              </LinkInternal>
            </ListItem>
            {isAdmin && (
              <ListItem>
                <LinkInternal
                  iconID={faArchive}
                  href="/account/administrator/audits"
                >
                  Audit Logs
                </LinkInternal>
              </ListItem>
            )}
            <ListItem>
              <LinkInternal iconID={faSignOutAlt} href="/authentication/logout">
                Logout
              </LinkInternal>
            </ListItem>
          </>
        )}
      </ListUnordered>
    </ListItem>
  );
}
