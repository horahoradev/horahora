import { faUser } from "@fortawesome/free-solid-svg-icons";

import { IProfileData } from "./types";

import { Icon } from "#components/icons";
import { Article, ArticleBody, ArticleHeader, IArticleProps } from "#components/article";
import { Heading } from "#components/headings";
import { blockComponent } from "#components/meta";

// eslint-disable-next-line
import styles from "./article.module.scss";

export interface IProfileArticleProps extends IArticleProps {
  profile: IProfileData;
}

export const ProfileArticle = blockComponent(styles.block, Component);

function Component({
  profile,
  headingLevel,
  ...blockProps
}: IProfileArticleProps) {
  const { Username, banned } = profile;

  return (
    <Article {...blockProps}>
      <ArticleHeader>
        <Heading>
          {Username} <b> {banned && "(banned)"} </b>
        </Heading>
        <div>
          <Icon icon={faUser} />
        </div>
      </ArticleHeader>
      <ArticleBody>
        <List></List>
      </ArticleBody>
    </Article>
  );
}
