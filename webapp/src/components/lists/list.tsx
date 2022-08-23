import clsx from "clsx";

// eslint-disable-next-line
import styles from "./list.module.scss";
import { blockComponent, type IBlockProps } from "#components/meta";

export type IListProps = (IListUnorderedProps | IListOrderedProps) & {
  type?: "ordered" | "unordered";
};
export interface IListOrderedProps extends IBlockProps<"ol"> {
  isHorizontal?: boolean;
}

export interface IListUnorderedProps extends IBlockProps<"ul"> {
  isHorizontal?: boolean;
}
export interface IListItemProps extends IBlockProps<"li"> {}

export const List = blockComponent(styles.block, Component);
export const ListOrdered = blockComponent(styles.ordered, OrderedComponent);
export const ListUnordered = blockComponent(
  styles.unordered,
  UnorderedComponent
);
export const ListItem = blockComponent(styles.item, ItemComponent);

function Component({
  type = "unordered",
  className,
  ...blockProps
}: IListProps) {
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

function OrderedComponent({
  children,
  className,
  isHorizontal = false,
  ...blockProps
}: IListOrderedProps) {
  const finalClassName = clsx(className, isHorizontal && styles.horizontal);

  return (
    <ol className={finalClassName} {...blockProps}>
      {children}
    </ol>
  );
}

function UnorderedComponent({
  children,
  className,
  isHorizontal = false,
  ...blockProps
}: IListUnorderedProps) {
  const finalClassName = clsx(className, isHorizontal && styles.horizontal);

  return (
    <ul className={finalClassName} {...blockProps}>
      {children}
    </ul>
  );
}

function ItemComponent({ children, ...blockProps }: IListItemProps) {
  return <li {...blockProps}>{children}</li>;
}
