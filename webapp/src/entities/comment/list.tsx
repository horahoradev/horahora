import { useEffect, useState } from "react";

import { CommentCard } from "./card";

import { CardList } from "#components/lists";
import { blockComponent, type IChildlessBlockProps } from "#components/meta";
import { IVideo } from "#entities/post";
import { getPostComments } from "#api/lib";
import { LoadingBar } from "#components/loading-bar";
import { NewCommentForm } from "#components/comments";

export interface ICommentListProps extends IChildlessBlockProps<"div"> {
  postID: IVideo["VideoID"];
}

export const CommentList = blockComponent(undefined, Component);

function Component({ postID, ...blockProps }: ICommentListProps) {
  const [comments, changeComments] =
    useState<Awaited<ReturnType<typeof getPostComments>>>();

  useEffect(() => {
    (async () => {
      const newComments = await getPostComments(postID);
      changeComments(newComments);
    })();
  }, [postID]);

  return (
    <div {...blockProps}>
      <CardList>
        {!comments ? (
          <LoadingBar />
        ) : (
          comments.map((comment) => (
            <CommentCard key={comment.id} comment={comment} />
          ))
        )}
      </CardList>
      <NewCommentForm
        videoID={postID}
        onNewComment={async (commentInit) => {
          const newComments = await getPostComments(postID);
          changeComments(newComments);
        }}
      />
    </div>
  );
}
