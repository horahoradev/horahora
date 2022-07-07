import styles from "./textarea.module.scss";

import { blockComponent, type IBlockProps } from "#components/meta";

export interface ITextAreaProps extends IBlockProps<"textarea"> {
  id: string;
  name: string;
}

export const TextArea = blockComponent(styles.block, Component);

export function Component({ children, ...blockProps }: ITextAreaProps) {
  return <textarea {...blockProps}>{children}</textarea>;
}
