import {
  faArchive,
  faBars,
  faKey,
  faSignOutAlt,
  faUser,
  faUpload,
} from "@fortawesome/free-solid-svg-icons";
import { useState, useEffect } from "react";

import styles from "./account-nav.module.scss";
import { ThemeSwitcher } from "./theme-switcher";

import { getUserdata } from "#api/index";
import { LinkInternal } from "#components/links";
import { UserRank } from "#api/types";
import { ListItem, ListUnordered } from "#components/lists";
import { Button } from "#components/buttons";
import { type ILoggedInUserData } from "#codegen/schema/001_interfaces";

export function AccountNavigation() {
  const [userData, setUserData] = useState<ILoggedInUserData>();
  const isRegistered = Boolean(userData && userData.username);
  const isAdmin = userData?.rank === UserRank.ADMIN;

  useEffect(() => {
    let ignore = false;

    (async () => {
      let userData = await getUserdata();
      if (!ignore) setUserData(userData);
    })();

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <ListUnordered className={styles.block}>
      <ListItem>
        <Button iconID={faBars}>Account</Button>
        <ListUnordered>
          {!isRegistered ? (
            <>
              <ListItem>
                <LinkInternal href="/authentication/login">Login</LinkInternal>
              </ListItem>
              <ListItem>
                <LinkInternal href="/authentication/register">
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
                <LinkInternal
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
    </ListUnordered>
  );
}
