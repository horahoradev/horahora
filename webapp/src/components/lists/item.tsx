import { blockComponent, type IBlockProps } from "#components/meta";

export interface IListItemProps extends IBlockProps<"li"> {}

export const ListItem = blockComponent(undefined, Component);

function Component({ children, ...blockProps }: IListItemProps) {
  return <li {...blockProps}>{children}</li>;
}
