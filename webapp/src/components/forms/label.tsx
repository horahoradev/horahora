import styles from "./label.module.scss";

import { blockComponent, type IBlockProps } from "#components/meta";
export interface ILabelProps extends IBlockProps<"label"> {
  htmlFor: string;
}

export const Label = blockComponent(styles.block, Component);

export function Component({ children, ...blockProps }: ILabelProps) {
  return <label {...blockProps}>{children}</label>;
}
