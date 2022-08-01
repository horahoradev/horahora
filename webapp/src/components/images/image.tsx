// eslint-disable-next-line
import { blockComponent, type IBlockProps } from "#components/meta";

// eslint-disable-next-line
import { ImageHTML, type IImageHTMLProps } from "./html";
// eslint-disable-next-line
import styles from "./image.module.scss";

export interface IImageProps extends IBlockProps<"span"> {
  src: IImageHTMLProps["src"];
}

export const Image = blockComponent(styles.block, Component);

function Component({ src, ...blockProps }: IImageProps) {
  return (
    <span {...blockProps}>
      <ImageHTML src={src} />
    </span>
  );
}
