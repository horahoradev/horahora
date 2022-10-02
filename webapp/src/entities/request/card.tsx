import { useEffect, useState } from "react";

import { type IArchivalRequest } from "./types";
import { ArchiveRequestStatus } from "./status";

import {
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  DL,
  DS,
  ListItem,
  ListUnordered,
  type ICardProps,
} from "#components/lists";
import { blockComponent } from "#components/meta";
import { Heading } from "#components/headings";
import { fetchProfile } from "#api/lib";
import { deleteArchivalRequest, retryArchivalRequest } from "#api/archives";
import { LinkInternal } from "#components/links";
import { LoadingBar } from "#components/loading-bar";
import { ProfileURL } from "#lib/urls";
import { Button } from "#components/buttons";

// eslint-disable-next-line
import styles from "./card.module.scss";

export interface IRequestCardProps extends ICardProps {
  request: IArchivalRequest;
  onUpdate: Function;
}

export const RequestCard = blockComponent(styles.block, Component);

function Component({
  request,
  onUpdate,
  headingLevel,
  ...blockProps
}: IRequestCardProps) {
  const [account, changeAccount] =
    useState<Awaited<ReturnType<typeof fetchProfile>>>();
  const {
    ArchivedVideos,
    BackoffFactor,
    CurrentTotalVideos,
    DownloadID,
    LastSynced,
    UndownloadableVideos,
    Url,
    UserID,
  } = request;

  useEffect(() => {
    (async () => {
      const newAccount = await fetchProfile(UserID, 1);
      changeAccount(newAccount);
    })();
  }, [UserID]);

  async function deleteRequest(download_id: number) {
    const formParams = new URLSearchParams([
      ["download_id", String(download_id)],
    ]);
    var resp = await deleteArchivalRequest(formParams);
    onUpdate();
    return resp;
  }

  async function retryRequest(download_id: number) {
    const formParams = new URLSearchParams([
      ["download_id", String(download_id)],
    ]);
    var resp = await retryArchivalRequest(formParams);
    onUpdate();
    return resp;
  }

  return (
    <Card id={`download-request-${DownloadID}`} {...blockProps}>
      <CardHeader>
        <Heading level={headingLevel}>{Url}</Heading>
      </CardHeader>
      <CardBody>
        <DL>
          <DS
            isHorizontal
            dKey="Status"
            dValue={<ArchiveRequestStatus request={request} />}
          />
          <DS
            isHorizontal
            dKey="Account"
            dValue={
              !account ? (
                <LoadingBar />
              ) : (
                <LinkInternal
                  href={new ProfileURL(account.UserID)}
                  target="_blank"
                >
                  {account.Username}
                </LinkInternal>
              )
            }
          />
          <DS isHorizontal dKey="Archived videos" dValue={ArchivedVideos} />
          <DS
            isHorizontal
            dKey="Current total videos"
            dValue={CurrentTotalVideos}
          />
          <DS
            isHorizontal
            dKey="Undownloadable videos"
            dValue={UndownloadableVideos}
          />
          <DS isHorizontal dKey="Backoff factor" dValue={BackoffFactor} />
          <DS dKey="Last Sync" dValue={LastSynced} />
        </DL>
      </CardBody>
      <CardFooter>
        <ListUnordered className={styles.actions} isHorizontal>
          <ListItem>
            <Button
              className={styles.delete}
              onClick={() => deleteRequest(DownloadID)}
            >
              Delete
            </Button>
          </ListItem>
          <ListItem>
            <LinkInternal href={"/account/archives/events/" + DownloadID}>Logs</LinkInternal>
          </ListItem>
          <ListItem>
            <Button
              className={styles.retry}
              onClick={() => retryRequest(DownloadID)}
            >
              Retry
            </Button>
          </ListItem>
        </ListUnordered>
      </CardFooter>
    </Card>
  );
}
