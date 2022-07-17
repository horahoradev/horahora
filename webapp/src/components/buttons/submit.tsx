import { type IBaseButtonProps, BaseButton } from "./base";

export interface IButtonSubmitProps extends Omit<IBaseButtonProps, "type"> {}

export function ButtonSubmit({ children, ...blockProps }: IButtonSubmitProps) {
  return (
    <BaseButton {...blockProps} type="submit">
      {children}
    </BaseButton>
  );
}
