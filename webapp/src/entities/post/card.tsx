import { type IVideo } from "#codegen/schema/001_interfaces";
import { Heading } from "#components/headings";
import { ImageLink } from "#components/images";
import { LinkInternal } from "#components/links";
import {
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  DL,
  DS,
  type ICardProps,
} from "#components/lists";
import { blockComponent } from "#components/meta";
import { ProfileURL } from "#lib/urls";

// eslint-disable-next-line
import styles from "./card.module.scss";

export interface IPostVideoCardProps extends ICardProps {
  post: IVideo;
}

/**
 * Not showing rating and views because they are not needed on the card list level.
 *
 * Bait title and a bait preview along with the author link
 * are more than enough to make a decision to go to the full page.
 */
export const PostCard = blockComponent(undefined, Component);

function Component({ post, headingLevel, ...blockProps }: IPostVideoCardProps) {
  const { Title, VideoID, ThumbnailLoc, AuthorID, AuthorName } = post;

  return (
    <Card {...blockProps}>
      <CardHeader>
        <a href={`/videos/${VideoID}`}>
          <Heading level={headingLevel}>{Title}</Heading>
        </a>
      </CardHeader>

      <CardFooter>
      <ImageLink
          className={styles.preview}
          src={ThumbnailLoc}
          href={`/videos/${VideoID}`}
          alt={Title}
          onError={(e) => {
            const img = e.target as HTMLImageElement;
            img.onerror = null;
            img.src = `${ThumbnailLoc.slice(0, -6)}.jpg`;
          }}
        />
        <DL>
          <DS
            isHorizontal
            dKey="By"
            dValue={<LinkInternal href={new ProfileURL(AuthorID)}> {AuthorName} </LinkInternal>}

              // @TODO: `AuthorID` should point to the imported profile
              // not the account which made the request.
              // <LinkInternal
              //   // @ts-expect-error schema type
              //   href={new ProfileURL(AuthorID)}
              //   rel={"referrer"}
              //   target={"_blank"}
              // >
              //   {AuthorName}
              // </LinkInternal>

          />
        </DL>
      </CardFooter>
    </Card>
  );
}
