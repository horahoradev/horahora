import { blockComponent, IBlockProps } from "#components/meta";

// eslint-disable-next-line
import styles from "./pre.module.scss";

export interface IPreProps extends IBlockProps<"pre"> {
  children?: string;
}

export const Pre = blockComponent(styles.block, Component);

function Component({ children, ...blockProps }: IPreProps) {
  return <pre {...blockProps}>{children}</pre>;
}
