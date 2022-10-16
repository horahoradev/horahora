import { blockComponent, type IChildlessBlockProps } from "#components/meta";

// eslint-disable-next-line
import styles from "./html.module.scss";

export interface IImageHTMLProps
  extends Omit<IChildlessBlockProps<"img">, "src"> {
  src: URL | string;
}

export const ImageHTML = blockComponent(styles.block, Component);

function Component({ src, alt = "", onError, ...blockProps }: IImageHTMLProps) {
  return (
    <img
      src={typeof src == "string" || src == null ? src : src.toString()}
      alt={alt}
      onError={onError ?? onError}
      {...blockProps}
    />
  );
}
