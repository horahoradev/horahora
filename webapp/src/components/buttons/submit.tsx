import { type IButtonHTMLProps, ButtonHTML } from "./html";

import { blockComponent } from "#components/meta";

export interface IButtonSubmitProps extends Omit<IButtonHTMLProps, "type"> {}

export const ButtonSubmit = blockComponent(undefined, Component);

export function Component({ children, ...blockProps }: IButtonSubmitProps) {
  return (
    <ButtonHTML {...blockProps} type="submit">
      {children}
    </ButtonHTML>
  );
}
