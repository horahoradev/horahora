import { type IArchivalEvent } from "./types";

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

// eslint-disable-next-line
import styles from "./card.module.scss";

export interface IEventCardProps extends ICardProps {
  event: IArchivalEvent;
}

export const EventCard = blockComponent(styles.block, Component);

function Component({ event, headingLevel, ...blockProps }: IEventCardProps) {
  const { message, parent_url, timestamp, video_url } = event;

  return (
    <Card {...blockProps}>
      <CardHeader>
        <Heading level={headingLevel}>{video_url}</Heading>
      </CardHeader>
      <CardBody className={styles.body}>{message}</CardBody>
      <CardFooter>
        <DL>
          <DS dKey="Date" dValue={timestamp} />
          <DS dKey="Parent" dValue={parent_url} />
        </DL>
      </CardFooter>
    </Card>
  );
}
