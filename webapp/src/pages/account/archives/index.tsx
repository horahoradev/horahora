import styles from "./index.module.scss";

import { Page } from "#components/page";
import { LinkInternal } from "#components/links";

function ArchivesPage() {
  return (
    <Page>
      <h1 className={styles.heading}>Archives</h1>
      <ul>
        <li>
          <LinkInternal href="/account/archives/requests">
            Requests
          </LinkInternal>
        </li>
        <li>
          <LinkInternal href="/account/archives/events">Events</LinkInternal>
        </li>
        <li>
          <LinkInternal href="/account/archives/downloads">
            Download Progress
          </LinkInternal>
        </li>
      </ul>
    </Page>
  );
}

export default ArchivesPage;
