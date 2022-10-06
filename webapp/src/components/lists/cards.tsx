import { type IHeadingLevel } from "#components/headings";
import { blockComponent, type IBlockProps } from "#components/meta";

// eslint-disable-next-line
import styles from "./cards.module.scss";

export interface ICardListProps extends IBlockProps<"div"> {
  /**
   * Is not used by the list itself,
   * but the components extending it have to incorporate it.
   */
  headingLevel?: IHeadingLevel;
}
export interface ICardProps extends IBlockProps<"article"> {
  /**
   * Is not used by the card itself,
   * but the components extending it have to incorporate it.
   */
  headingLevel?: IHeadingLevel;
}
export interface ICardHeaderProps extends IBlockProps<"header"> {}
export interface ICardBodyProps extends IBlockProps<"section"> {}
export interface ICardFooterProps extends IBlockProps<"footer"> {}

export const CardList = blockComponent(styles.block, CardListComponent);
export const StartCardList = blockComponent(styles.cardlist, CardListComponent);

/**
 * Skeleton component, do not use it directly but instead
 * create a card component which extends its props.
 */
export const Card = blockComponent(styles.card, CardComponent);
export const CardHeader = blockComponent(styles.header, CardHeaderComponent);
export const CardBody = blockComponent(styles.body, CardBodyComponent);
export const CardFooter = blockComponent(styles.footer, CardFooterComponent);

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
