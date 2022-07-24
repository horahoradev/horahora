import { type IButtonHTMLProps, ButtonHTML } from "./html";

import { blockComponent } from "#components/meta";

export interface IButtonProps extends Omit<IButtonHTMLProps, "type"> {}

export const Button = blockComponent(undefined, Component);

export function Component({ children, ...blockProps }: IButtonProps) {
  return (
    <ButtonHTML {...blockProps} type="button">
      {children}
    </ButtonHTML>
  );
}
