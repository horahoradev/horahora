import { useRouter } from "next/router";
import { useCallback } from "react";

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

// eslint-disable-next-line
import styles from "./article.module.scss";

export interface IPostArticleProps extends IArticleProps {
  video: IVideoDetail;
}

export const PostArticle = blockComponent(styles.block, Component);

function Component({ video, headingLevel, ...blockProps }: IPostArticleProps) {
  const router = useRouter();
  const {
    Title,
    MPDLoc,
    Views,
    RecommendedVideos,
    UploadDate,
    VideoID,
    Rating,
    Tags,
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
        <Heading level={headingLevel}>{Title}</Heading>
        <VideoPlayer url={MPDLoc} next_video={nextVideo} />
      </ArticleHeader>
      <ArticleBody>
        <DL>
          <DS
            dKey={"Tags"}
            dValue={
              !Tags.length ? (
                "None"
              ) : (
                <List className={styles.tags} isHorizontal>
                  {Tags.map((tag, idx) => {
                    return (
                      <ListItem key={idx}>
                        <LinkInternal href={`/?search=${tag}`} target="_blank">
                          {tag}
                        </LinkInternal>
                      </ListItem>
                    );
                  })}
                </List>
              )
            }
          />
          <DS isHorizontal dKey={"Views"} dValue={Views} />
          <DS isHorizontal dKey={"Upload date"} dValue={UploadDate} />
        </DL>
        <PostRate postID={VideoID} rating={Rating} />
      </ArticleBody>
      <ArticleFooter></ArticleFooter>
    </Article>
  );
}
