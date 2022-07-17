import { blockComponent, IBlockProps } from "#components/meta";

const headingLevels = [1, 2, 3, 4, 5, 6] as const;
type IHeadingLevel = typeof headingLevels[number];

export interface ICardListProps extends IBlockProps<"div"> {}
export interface ICardProps extends IBlockProps<"article"> {
  /**
   * Is not used by the card itself,
   * but the components extending it have incorporate it.
   */
  headingLevel: IHeadingLevel;
}
export interface ICardHeaderProps extends IBlockProps<"header"> {}
export interface ICardBodyProps extends IBlockProps<"section"> {}
export interface ICardFooterProps extends IBlockProps<"footer"> {}

export const CardList = blockComponent(undefined, CardListComponent);
/**
 * Skeleton component, do not use it directly but instead
 * create a card component which extends its props.
 */
export const Card = blockComponent(undefined, CardComponent);
export const CardHeader = blockComponent(undefined, CardHeaderComponent);
export const CardBody = blockComponent(undefined, CardBodyComponent);
export const CardFooter = blockComponent(undefined, CardFooterComponent);

function CardListComponent({ children, ...blockProps }: ICardListProps) {
  return <div {...blockProps}>{children}</div>;
}

function CardComponent({ children, ...blockProps }: ICardProps) {
  return <article {...blockProps}>{children}</article>;
}

function CardHeaderComponent({ children, ...blockProps }: ICardHeaderProps) {
  return <header {...blockProps}>{children}</header>;
}

function CardBodyComponent({ children, ...blockProps }: ICardBodyProps) {
  return <section {...blockProps}>{children}</section>;
}

function CardFooterComponent({ children, ...blockProps }: ICardFooterProps) {
  return <footer {...blockProps}>{children}</footer>;
}
