import { blockComponent, type IBlockProps } from "#components/meta";
import { type IHeadingLevel } from "#components/headings";

export interface IArticleProps extends IBlockProps<"article"> {
  /**
   * Is not used by the card itself,
   * but the components extending it have incorporate it.
   */
  headingLevel: IHeadingLevel;
}
export interface IArticleHeaderProps extends IBlockProps<"header"> {}
export interface IArticleBodyProps extends IBlockProps<"section"> {}
export interface IArticleFooterProps extends IBlockProps<"footer"> {}

/**
 * Skeleton component, do not use it directly but instead
 * create an article component which extends its props.
 * The main difference between this and card is card is a part of a list.
 * But article is one of a kind per page.
 */
export const Article = blockComponent(undefined, ArticleComponent);
export const ArticleHeader = blockComponent(undefined, ArticleHeaderComponent);
export const ArticleBody = blockComponent(undefined, ArticleBodyComponent);
export const ArticleFooter = blockComponent(undefined, ArticleFooterComponent);

function ArticleComponent({ children, ...blockProps }: IArticleProps) {
  return <article {...blockProps}>{children}</article>;
}
function ArticleHeaderComponent({
  children,
  ...blockProps
}: IArticleHeaderProps) {
  return <header {...blockProps}>{children}</header>;
}
function ArticleBodyComponent({ children, ...blockProps }: IArticleBodyProps) {
  return <section {...blockProps}>{children}</section>;
}
function ArticleFooterComponent({
  children,
  ...blockProps
}: IArticleFooterProps) {
  return <footer {...blockProps}>{children}</footer>;
}
