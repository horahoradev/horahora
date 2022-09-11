import {
  faCheck,
  faRotate,
  faBoxArchive,
} from "@fortawesome/free-solid-svg-icons";

import { type IArchivalRequest } from "#codegen/schema/001_interfaces";
import { Icon } from "#components/icons";

// eslint-disable-next-line
import styles from "./status.module.scss";

export interface IArchiverRequestStatusProps {
  request: IArchivalRequest;
}

export function ArchiveRequestStatus({ request }: IArchiverRequestStatusProps) {
  if (
    request.ArchivedVideos === request.CurrentTotalVideos &&
    request.CurrentTotalVideos != 0
  ) {
    return (
      <span className={styles.complete}>
        <Icon icon={faCheck} /> <span>Complete</span>
      </span>
    );
  } else if (request.CurrentTotalVideos == 0 || request.LastSynced == null) {
    return (
      <span className={styles.fetch}>
        <Icon icon={faRotate} /> <span>Fetching...</span>
      </span>
    );
  } else {
    return (
      <span className={styles.archive}>
        <Icon icon={faBoxArchive} /> <span>Archiving...</span>
      </span>
    );
  }
}
