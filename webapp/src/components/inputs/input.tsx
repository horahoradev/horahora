import { forwardRef, type LegacyRef } from "react";

import { blockComponent, type IChildlessBlockProps } from "#components/meta";

// eslint-disable-next-line
import styles from "./input.module.scss";

export interface IInputProps extends IChildlessBlockProps<"input"> {
  id: string;
  name: string;
}

export const Input = forwardRef<HTMLInputElement, IInputProps>(
  blockComponent(styles.block, Component)
);

function Component(
  { ...blockProps }: IInputProps,
  ref?: LegacyRef<HTMLInputElement>
) {
  return <input {...blockProps} ref={ref} />;
}
