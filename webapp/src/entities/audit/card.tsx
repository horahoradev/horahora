import { useEffect, useState } from "react";

import { type IAuditEvent } from "./types";

import {
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  DL,
  DS,
  type ICardProps,
} from "#components/lists";
import { blockComponent } from "#components/meta";
import { Heading } from "#components/headings";
import { LoadingBar } from "#components/loading-bar";
import { LinkInternal } from "#components/links";
import { ProfileURL } from "#lib/urls";
import { fetchProfile } from "#api/lib";

export interface IAuditCardProps extends Omit<ICardProps, "id"> {
  audit: IAuditEvent;
}

export const AuditCard = blockComponent(undefined, Component);

function Component({ audit, headingLevel, ...blockProps }: IAuditCardProps) {
  const [account, changeAccount] =
    useState<Awaited<ReturnType<typeof fetchProfile>>>();
  const { ID, Message, Timestamp, UserID } = audit;

  useEffect(() => {
    (async () => {
      const newAccount = await fetchProfile(UserID, 1);
      changeAccount(newAccount);
    })();
  }, [UserID]);

  return (
    <Card id={`audit-${ID}`} {...blockProps}>
      <CardHeader>
        <Heading level={headingLevel}>
          {!account ? (
            <LoadingBar />
          ) : (
            <LinkInternal href={new ProfileURL(account.UserID)} target="_blank">
              {account.Username}
            </LinkInternal>
          )}
        </Heading>
      </CardHeader>
      <CardBody>{Message}</CardBody>
      <CardFooter>
        <DL>
          <DS isHorizontal dKey="Date" dValue={Timestamp} />
        </DL>
      </CardFooter>
    </Card>
  );
}
