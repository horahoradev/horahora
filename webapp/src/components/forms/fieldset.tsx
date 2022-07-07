import { blockComponent, IBlockProps } from "#components/meta";

export interface IFieldsetProps extends IBlockProps<"fieldset"> {}

export const Fieldset = blockComponent(undefined, Component);

function Component({ children, ...blockProps }: IFieldsetProps) {
  return <fieldset {...blockProps}>{children}</fieldset>;
}
