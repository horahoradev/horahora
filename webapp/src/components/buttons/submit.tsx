import { type IButtonHTMLProps, ButtonHTML } from "./html";

export interface IButtonSubmitProps extends Omit<IButtonHTMLProps, "type"> {}

export function ButtonSubmit({ children, ...blockProps }: IButtonSubmitProps) {
  return (
    <ButtonHTML {...blockProps} type="submit">
      {children}
    </ButtonHTML>
  );
}
