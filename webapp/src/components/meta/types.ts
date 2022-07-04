import type { ComponentPropsWithoutRef } from "react";

/**
 * Base props for block components.
 */
export type IBlockProps<HTMLTag extends keyof JSX.IntrinsicElements> =
  ComponentPropsWithoutRef<HTMLTag>;

export type IFuncComponent<PropsType> = (
  props: PropsType,
  ref?: any
) => JSX.Element;
