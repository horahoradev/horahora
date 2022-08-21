import { IVideo } from "./types";

import { blockComponent, IBlockProps } from "#components/meta";
import { Button } from "#components/buttons";
import { ratePost } from "#api/lib";

// eslint-disable-next-line
import styles from "./rating.module.scss";

export interface IPostRateProps extends IBlockProps<"div"> {
  postID: IVideo["VideoID"]
}

export const PostRate = blockComponent(styles.block, Component);

function Component({ postID, ...blockProps }: IPostRateProps) {

  async function rate(rating: number) {
    if (id == 0) {
      // TODO: throw
      return;
    }
    await ratePost(id, rating);
    setRating(rating);
  }

  async function likePost() {}
  async function dislikePost() {}

  return (
    <div {...blockProps}>
      <Button>Dislike</Button>
      <Button>Like</Button>
    </div>
  );
}
