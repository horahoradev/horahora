import type { ComponentPropsWithoutRef } from "react";

/**
 * Base props for block components.
 */
export type IBlockProps<HTMLTag extends keyof JSX.IntrinsicElements> =
  ComponentPropsWithoutRef<HTMLTag>;

/**
 * Base block props but for empty components.
 */
export type IChildlessBlockProps<HTMLTag extends keyof JSX.IntrinsicElements> =
  Omit<IBlockProps<HTMLTag>, "children">;

export type IFuncComponent<PropsType> = (
  props: PropsType,
  ref?: any
) => JSX.Element;
