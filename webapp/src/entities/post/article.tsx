import { useRouter } from "next/router";
import { useCallback } from "react";
import { faUserCircle } from "@fortawesome/free-solid-svg-icons";

import { type IVideoDetail } from "./types";
import { PostRate } from "./rating";

import {
  Article,
  ArticleBody,
  ArticleFooter,
  ArticleHeader,
  type IArticleProps,
} from "#components/article";
import { Heading } from "#components/headings";
import { blockComponent } from "#components/meta";
import { VideoPlayer } from "#components/video";
import { DL, DS, List, ListItem } from "#components/lists";
import { LinkInternal } from "#components/links";
import { ProfileURL } from "#lib/urls";
import { VideoAdminControls } from "#components/account";
import { useAccount } from "#hooks";
import { CommentList } from "#entities/comment";

// eslint-disable-next-line
import styles from "./article.module.scss";
import { Pre } from "#components/pre";

export interface IPostArticleProps extends IArticleProps {
  video: IVideoDetail;
}

export const PostArticle = blockComponent(styles.block, Component);

function Component({
  video,
  headingLevel = 2,
  ...blockProps
}: IPostArticleProps) {
  const router = useRouter();
  const { isAdmin } = useAccount();
  const {
    Title,
    MPDLoc,
    Views,
    RecommendedVideos,
    UploadDate,
    VideoID,
    Rating,
    Tags,
    AuthorID,
    Username,
    VideoDescription,
  } = video;

  const nextVideo = useCallback(async () => {
    if (!RecommendedVideos.length) {
      return;
    }

    const [nextVid] = RecommendedVideos;
    router.push(`/videos/${nextVid.VideoID}`);
  }, [RecommendedVideos]);

  return (
    <Article {...blockProps}>
      <ArticleHeader>
        <Heading className={styles.title} level={headingLevel}>
          {Title}
        </Heading>
        <VideoPlayer url={MPDLoc} next_video={nextVideo} />
      </ArticleHeader>
      <PostRate postID={VideoID} rating={Rating} />
      <ArticleBody>
        <DL>
          <DS isHorizontal dKey="Views" dValue={Views} />
          <DS
            isHorizontal
            dKey="Author"
            dValue={
              <LinkInternal
                href={new ProfileURL(AuthorID)}
                iconID={faUserCircle}
              >
                {Username}
              </LinkInternal>
            }
          />
          <DS isHorizontal dKey="Upload date" dValue={UploadDate} />
          {!Tags ? null : (
            <DS
              dKey="Tags"
              dValue={
                <List className={styles.tags} isHorizontal>
                  {Tags.map((tag, idx) => {
                    return (
                      <ListItem className={styles.tag} key={idx}>
                        <LinkInternal href={`/?search=${tag}`} target="_blank">
                          {tag}
                        </LinkInternal>
                      </ListItem>
                    );
                  })}
                </List>
              }
            />
          )}
          <DS
            dKey="Description"
            dValue={
              <p>
                <Pre>{VideoDescription}</Pre>
              </p>
            }
          />
        </DL>
        {isAdmin && <VideoAdminControls data={video} />}
      </ArticleBody>
      <ArticleFooter>
        <Heading
          // @ts-expect-error numbers
          level={headingLevel + 1}
        >
          Comments
        </Heading>
        <CommentList postID={VideoID} />
      </ArticleFooter>
    </Article>
  );
}
