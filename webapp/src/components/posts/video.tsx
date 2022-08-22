import { Tag, Avatar, Rate, Comment, List } from "antd";
import { faThumbsUp, faUserCircle } from "@fortawesome/free-solid-svg-icons";
import { UserOutlined } from "@ant-design/icons";

import { ProfileURL } from "#lib/urls";
import { UserRank } from "#lib/account";
import { Icon } from "#components/icons";
import { upvoteComment as apiUpvoteComment, ratePost } from "#api/lib";
import { LinkInternal } from "#components/links";

import { NewCommentForm } from "#components/comments";
import { type IComment } from "#types/entities";
import { IVideoDetail } from "#entities/post";

// eslint-disable-next-line
import styles from "./video.module.scss";

export interface IVideoViewProps {
  id: number;
  data: IVideoDetail;
  videoComments: IComment[];

  refreshComments: () => Promise<unknown>;
}

export function VideoView(props: IVideoViewProps) {
  let {data, id, videoComments, refreshComments } = props;

  async function upvoteComment(commentID: number, has_upvoted: boolean) {
    await apiUpvoteComment(commentID, !has_upvoted);
    await refreshComments();
  }

  // FIXME: new API endpoint
  return (
    <div className={styles.block}>
      <List
        bordered={false}
        split={false}
        className="comment-list"
        header={
          <h2 className="ml-4 mb-0 text-xl text-black dark:text-white">
            Comments ({videoComments.length})
          </h2>
        }
        itemLayout="horizontal"
        dataSource={videoComments}
        renderItem={(item) => (
          <li>
            <Comment
              className="border-0 text-black dark:text-white shadow-none"
              actions={[
                <span key={1} className="text-black dark:text-white text-bold">
                  {item.upvote_count}
                </span>,
                <Icon
                  key={2}
                  onClick={() => upvoteComment(item.id, item.user_has_upvoted)}
                  className={
                    item.user_has_upvoted
                      ? "max-h-4 mr-1 text-green-400"
                      : "max-h-4 mr-1 text-gray-400"
                  }
                  icon={faThumbsUp}
                />,
              ]}
              author={
                <b className="text-black dark:text-white">{item.fullname}</b>
              }
              avatar={
                <Avatar
                  className="ml-4"
                  shape="square"
                  size={50}
                  icon={<UserOutlined />}
                />
              }
              content={item.content}
              datetime={item.created}
            />
          </li>
        )}
      />
      <NewCommentForm
        videoID={id}
        onNewComment={async (commentInit) => {
          await refreshComments();
        }}
      />
    </div>
  );
}
