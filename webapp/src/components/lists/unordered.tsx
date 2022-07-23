import { blockComponent, type IBlockProps } from "#components/meta";

export interface IListUnorderedProps extends IBlockProps<"ul"> {}

export const ListUnordered = blockComponent(undefined, Component);

function Component({ children, ...blockProps }: IListUnorderedProps) {
  return <ul {...blockProps}>{children}</ul>;
}
