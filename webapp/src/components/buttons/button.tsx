import { type IButtonHTMLProps, ButtonHTML } from "./html";

export interface IButtonProps extends Omit<IButtonHTMLProps, "type"> {}
export function Button({ children, ...blockProps }: IButtonProps) {
  return (
    <ButtonHTML {...blockProps} type="button">
      {children}
    </ButtonHTML>
  );
}
