import { type IBaseButtonProps, BaseButton } from "./base";

export interface IButtonProps extends Omit<IBaseButtonProps, "type"> {}
export function Button({ children, ...blockProps }: IButtonProps) {
  return (
    <BaseButton {...blockProps} type="button">
      {children}
    </BaseButton>
  );
}
