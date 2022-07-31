import { ImageHTML, type IImageHTMLProps } from "#components/images";
import { LinkExternal, type ILinkExternalProps } from "#components/links";
import { blockComponent } from "#components/meta";

// eslint-disable-next-line
import styles from "./link.module.scss";

export interface IImageLinkProps extends ILinkExternalProps {
  src: IImageHTMLProps["src"];
  alt: IImageHTMLProps["alt"];

  onImageError: IImageHTMLProps["onError"];
}

/**
 * Image but also a clickable link
 */
export const ImageLink = blockComponent(styles.block, Component);

function Component({
  src,
  alt,
  href,
  onImageError,
  ...blockProps
}: IImageLinkProps) {
  return (
    <LinkExternal href={href ?? src} {...blockProps}>
      <ImageHTML src={src} alt={alt} onError={onImageError ?? onImageError} />
    </LinkExternal>
  );
}
