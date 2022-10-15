import { PostCard } from "./card";

import { blockComponent } from "#components/meta";
import { type IVideo } from "#codegen/schema/001_interfaces";
import { CardList, type ICardListProps } from "#components/lists";
import { LinkInternal } from "#components/links";

export interface IPostListProps extends ICardListProps {
  posts: IVideo[];
}

export const PostList = blockComponent(undefined, Component);

function Component({
  posts = [],
  headingLevel,
  ...blockProps
}: IPostListProps) {
  return (
    <CardList {...blockProps}>
      {posts.map((post) => (
        <PostCard key={post.VideoID} headingLevel={headingLevel} post={post} {...<LinkInternal href="/"></LinkInternal>}></PostCard>
      ))}
    </CardList>
  );
}
