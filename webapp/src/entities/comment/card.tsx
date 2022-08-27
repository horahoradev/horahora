import { useCallback, useState } from "react";
import { faThumbsUp } from "@fortawesome/free-solid-svg-icons";
import clsx from "clsx";

import { type ICommentData } from "./types";

import { Heading } from "#components/headings";
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
import { Button } from "#components/buttons";
import { upvoteComment } from "#api/lib";

// eslint-disable-next-line
import styles from "./card.module.scss";

export interface ICommentCardProps extends Omit<ICardProps, "id"> {
  comment: ICommentData;
}

export const CommentCard = blockComponent(styles.block, Component);

function Component({
  comment,
  headingLevel,
  ...blockProps
}: ICommentCardProps) {
  const [upvoteCount, changeUpvoteCount] = useState(comment.upvote_count);
  const [isUpvoted, changeIsUpvoted] = useState(comment.user_has_upvoted);
  const { id, fullname, created, content } = comment;

  const upvote = useCallback(async () => {
    await upvoteComment(id, isUpvoted);

    changeUpvoteCount((upvoteCount) =>
      isUpvoted ? upvoteCount - 1 : upvoteCount + 1
    );
    changeIsUpvoted(!isUpvoted);
  }, [id, isUpvoted]);

  return (
    <Card id={`comment-${id}`} {...blockProps}>
      <CardHeader>
        <Heading level={headingLevel}>{fullname}</Heading>
      </CardHeader>
      <CardBody>{content}</CardBody>
      <CardFooter>
        <DL>
          <DS
            isHorizontal
            dKey="Upvotes"
            dValue={
              <Button
                className={clsx(isUpvoted && styles.upvoted)}
                iconID={faThumbsUp}
                onClick={async () => {
                  await upvote();
                }}
              >
                {upvoteCount}
              </Button>
            }
          />
          <DS isHorizontal dKey="Date created" dValue={created} />
        </DL>
      </CardFooter>
    </Card>
  );
}
