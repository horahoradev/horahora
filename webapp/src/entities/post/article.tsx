import { useRouter } from "next/router";
import { useCallback } from "react";

import { type IVideoDetail } from "./types";

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

// eslint-disable-next-line
import styles from "./article.module.scss";
import { DL, DS } from "#components/lists";

export interface IPostArticleProps extends IArticleProps {
  video: IVideoDetail;
}

export const PostArticle = blockComponent(styles.block, Component);

function Component({ video, headingLevel, ...blockProps }: IPostArticleProps) {
  const router = useRouter();
  const { Title, MPDLoc, Views, RecommendedVideos, UploadDate } = video;

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
      </ArticleHeader>
      <ArticleBody>
        <VideoPlayer
          className={styles.player}
          url={MPDLoc}
          next_video={nextVideo}
        />
        <DL>
          <DS dKey={"Views"} dValue={Views} />
          <DS dKey={"Upload date"} dValue={UploadDate} />
        </DL>
      </ArticleBody>
      <ArticleFooter></ArticleFooter>
    </Article>
  );
}
