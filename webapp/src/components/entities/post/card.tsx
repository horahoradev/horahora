import { type IVideo } from "#codegen/schema/001_interfaces";
import { Heading } from "#components/headings";
import { LinkInternal } from "#components/links";
import {
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  type ICardProps,
} from "#components/lists";
import { blockComponent } from "#components/meta";

export interface IPostVideoCardProps extends ICardProps {
  post: IVideo;
}

/**
 * Not showing rating and views because they are not needed on the card list level.
 *
 * Bait title and a bait preview are more than enough
 * to make a decision to go to the full page.
 */
export const PostCard = blockComponent(undefined, Component);

function Component({
  post,
  headingLevel = 2,
  ...blockProps
}: IPostVideoCardProps) {
  const { Title, VideoID, ThumbnailLoc } = post;

  return (
    <Card {...blockProps}>
      <CardHeader>
        <Heading level={headingLevel}>{Title}</Heading>
      </CardHeader>
      <CardBody>
        <img
          alt={Title}
          src={`${ThumbnailLoc}`}
          onError={(e) => {
            const img = e.target as HTMLImageElement;
            img.onerror = null;
            // @ts-expect-error schema type
            img.src = `${ThumbnailLoc.slice(0, -6)}.jpg`;
          }}
        />
      </CardBody>
      <CardFooter>
        <LinkInternal href={`/videos/${VideoID}`}>Details</LinkInternal>
      </CardFooter>
    </Card>
  );
}
