import styles from "./input.module.scss";

import { blockComponent, type IChildlessBlockProps } from "#components/meta";

export interface IInputProps extends IChildlessBlockProps<"input"> {
  id: string,
  name: string
}

export const Input = blockComponent(styles.block, Component);

function Component({ ...blockProps }: IInputProps) {
  return <input {...blockProps} />;
}
