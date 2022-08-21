import { faUser } from "@fortawesome/free-solid-svg-icons";

import { IProfileData } from "./types";

import { Icon } from "#components/icons";
import {
  Article,
  ArticleBody,
  ArticleFooter,
  ArticleHeader,
  IArticleProps,
} from "#components/article";
import { Heading } from "#components/headings";
import { blockComponent } from "#components/meta";
import {
  banAccount,
  promoteAccountToMod,
  promoteAccountToAdmin,
} from "#api/lib";
import { UserRank } from "#lib/account";
import { Button } from "#components/buttons";

// eslint-disable-next-line
import styles from "./article.module.scss";
import { List } from "#components/lists";

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
      <ArticleHeader className={styles.header}>
        <Heading>
          {Username} <b> {banned && "(banned)"} </b>
        </Heading>
        <Icon icon={faUser} />
      </ArticleHeader>
      <ArticleBody></ArticleBody>
      <ArticleFooter>
        <List></List>
      </ArticleFooter>

      {profile.L && profile.L.rank === UserRank.ADMIN && (
        <p className={"flex justify-center"}>
          <Button onClick={() => banAccount(profile.UserID)}>Ban</Button>
        </p>
      )}
      {profile.L && profile.L.rank === UserRank.ADMIN && (
        <p className={"flex justify-center"}>
          <Button onClick={() => promoteAccountToMod(profile.UserID)}>
            Promote to mod
          </Button>
        </p>
      )}
      {profile.L && profile.L.rank === UserRank.ADMIN && (
        <p className={"flex justify-center"}>
          <Button onClick={() => promoteAccountToAdmin(profile.UserID)}>
            Promote to admin
          </Button>
        </p>
      )}

      {/* TODO: no more copy pasta! */}
    </Article>
  );
}
