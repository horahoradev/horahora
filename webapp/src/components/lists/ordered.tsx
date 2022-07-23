import { blockComponent, type IBlockProps } from "#components/meta";

export interface IListOrderedProps extends IBlockProps<"ol"> {}

export const ListOrdered = blockComponent(undefined, Component);

function Component({ children, ...blockProps }: IListOrderedProps) {
  return <ol {...blockProps}>{children}</ol>;
}
