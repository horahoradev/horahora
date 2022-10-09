import { Input, type IInputProps } from "./input";

import { blockComponent } from "#components/meta";

export interface IHiddenProps extends Omit<IInputProps, "type"> {}

export const Hidden = blockComponent(undefined, Component);

function Component({ ...blockProps }: IHiddenProps) {
  return <Input {...blockProps} type="hidden" />;
}
