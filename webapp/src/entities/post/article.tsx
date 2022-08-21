import { type IVideoDetail } from "./types";

import {
  Article,
  ArticleHeader,
  type IArticleProps,
} from "#components/article";
import { Heading } from "#components/headings";
import { blockComponent } from "#components/meta";

// eslint-disable-next-line
import styles from "./article.module.scss";

export interface IPostArticleProps extends IArticleProps {
  video: IVideoDetail;
}

export const PostArticle = blockComponent(styles.block, Component);

function Component({ video, headingLevel, ...blockProps }: IPostArticleProps) {
  const {} = video;

  return (
    <Article {...blockProps}>
      <ArticleHeader>
        <Heading level={headingLevel}></Heading>
      </ArticleHeader>
    </Article>
  );
}
