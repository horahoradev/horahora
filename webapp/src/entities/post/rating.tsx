import { useState } from "react";
import { faThumbsDown, faThumbsUp } from "@fortawesome/free-solid-svg-icons";
import clsx from "clsx";

import { IVideo } from "./types";

import { blockComponent, IBlockProps } from "#components/meta";
import { Button } from "#components/buttons";
import { ratePost } from "#api/lib";

// eslint-disable-next-line
import styles from "./rating.module.scss";

export interface IPostRateProps extends IBlockProps<"div"> {
  postID: IVideo["VideoID"];
  rating?: IVideo["Rating"];
}

export const PostRate = blockComponent(styles.block, Component);

function Component({
  postID,
  rating,
  className,
  ...blockProps
}: IPostRateProps) {
  const [isLiked, changeLike] = useState(
    rating === undefined ? undefined : rating === 5
  );
  const finalClassName = clsx(
    className,
    isLiked === undefined
      ? undefined
      : isLiked
      ? styles.block_liked
      : styles.block_disliked
  );

  async function likePost() {
    await ratePost(postID, 5);
    changeLike(true);
  }
  async function dislikePost() {
    await ratePost(postID, 0);
    changeLike(false);
  }

  return (
    <div className={finalClassName} {...blockProps}>
      <Button
        className={styles.dislike}
        iconID={faThumbsDown}
        onClick={async () => {
          await dislikePost();
        }}
      >
        Dislike
      </Button>
      <Button
        className={styles.like}
        iconID={faThumbsUp}
        onClick={async () => {
          await likePost();
        }}
      >
        Like
      </Button>
    </div>
  );
}
