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
import { List, ListItem } from "#components/lists";
import { useAccount } from "#hooks";

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
  const { isAdmin } = useAccount();
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

      {isAdmin && (
        <ArticleFooter>
          <List isHorizontal>
            <ListItem>
              <Button onClick={() => banAccount(profile.UserID)}>Ban</Button>
            </ListItem>
            <ListItem>
              <Button onClick={() => promoteAccountToMod(profile.UserID)}>
                Promote to mod
              </Button>
            </ListItem>
            <ListItem>
              <Button onClick={() => promoteAccountToAdmin(profile.UserID)}>
                Promote to admin
              </Button>
            </ListItem>
          </List>
        </ArticleFooter>
      )}
    </Article>
  );
}
