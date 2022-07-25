import styles from "./header.module.scss";
import { Search } from "./search";
import { UserNav } from "./account";

import { LinkInternal } from "#components/links/internal";

interface IHeaderProps {
  userData?: Record<string, unknown>;
  dataless?: boolean;
}

export function Header({ userData, dataless }: IHeaderProps) {
  return (
    <header className={styles.block}>
      <nav className="max-w-screen-lg w-screen flex justify-start items-center gap-x-4 mx-4">
        <div className="flex justify-start flex-grow-0">
          {/* @TODO: site logo component */}
          <LinkInternal className={styles.logo} href="/">
            Horahora
          </LinkInternal>
        </div>
        <Search />
        {!dataless && (
          <div className="flex-grow-0 ml-auto">
            <UserNav userData={userData} />
          </div>
        )}
      </nav>
    </header>
  );
}
