import { useEffect, useState } from "react";
import clsx from "clsx";
import type { ComponentPropsWithoutRef, ReactElement } from "react";

import { IS_BROWSER } from "#environment/constants";

/**
 * Base props for block components.
 */
export type IBlockProps<T extends keyof JSX.IntrinsicElements> =
  ComponentPropsWithoutRef<T>;

type FuncComponent<PropsType> = (props: PropsType, ref?: any) => JSX.Element;

/**
 * A decorator for components with a root element.
 * `className` prop passed to it will be appended to the `classList`
 * instead of overwriting the class attribute.
 * @param blockClassName `className` of the root element
 * @param functionComponent
 * @returns Wrapped function.
 */
export function blockComponent<Props>(
  blockClassName: string | string[] | undefined,
  functionComponent: FuncComponent<Props>
): FuncComponent<Props> {
  // @ts-expect-error
  return ({ className, ...blockProps }: Props, ref) => {
    const baseClass = clsx(blockClassName, className);

    // @ts-expect-error
    return functionComponent({ className: baseClass, ...blockProps }, ref);
  };
}

export interface IClientComponentProps {
  children: ReactElement;
}

/**
 * Renders the children only on client.
 */
export function ClientComponent({ children }: IClientComponentProps) {
  const [isEnabled, enableComponent] = useState(false);

  useEffect(() => {
    IS_BROWSER && enableComponent(true);
  }, []);

  return <>{isEnabled && children}</>;
}
