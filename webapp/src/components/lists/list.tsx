import { blockComponent, type IBlockProps } from "#components/meta";

export type IListProps = (IListUnorderedProps | IListOrderedProps) & {
  type?: "ordered" | "unordered";
};
export interface IListOrderedProps extends IBlockProps<"ol"> {}

export interface IListUnorderedProps extends IBlockProps<"ul"> {}
export interface IListItemProps extends IBlockProps<"li"> {}

export const List = blockComponent(undefined, Component);
export const ListOrdered = blockComponent(undefined, OrderedComponent);
export const ListUnordered = blockComponent(undefined, UnorderedComponent);
export const ListItem = blockComponent(undefined, ItemComponent);

function Component({ type = "unordered", ...blockProps }: IListProps) {
  switch (type) {
    case "unordered": {
      return <ListUnordered {...blockProps} />;
    }
    case "ordered": {
      return <ListOrdered {...blockProps} />;
    }

    default: {
      throw new Error(`Illegal list type of "${type}".`);
    }
  }
}

function OrderedComponent({ children, ...blockProps }: IListOrderedProps) {
  return <ol {...blockProps}>{children}</ol>;
}

function UnorderedComponent({ children, ...blockProps }: IListUnorderedProps) {
  return <ul {...blockProps}>{children}</ul>;
}

function ItemComponent({ children, ...blockProps }: IListItemProps) {
  return <li {...blockProps}>{children}</li>;
}
