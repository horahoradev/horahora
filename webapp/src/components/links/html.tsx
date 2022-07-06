import { forwardRef, type LegacyRef } from "react";
import { type IconDefinition } from "@fortawesome/fontawesome-svg-core";

import styles from "./html.module.scss";

import { blockComponent, IBlockProps } from "#components/meta";
import { Icon } from "#components/icons";
export interface IAnchourHTMLProps extends Omit<IBlockProps<"a">, "href"> {
  /**
   * Icon id to use in a link.
   */
  iconID?: IconDefinition;
  href?: string | URL;
}

/**
 * Baseline link component for other link components.
 */
// needs to be coated in `forwardRef()` to work with nextjs `Link` element.
export const AnchourHTML = forwardRef<HTMLAnchorElement, IAnchourHTMLProps>(
  blockComponent(styles.block, Component)
);

function Component(
  { href, iconID, children, ...blockProps }: IAnchourHTMLProps,
  ref?: LegacyRef<HTMLAnchorElement>
) {
  return (
    <a
      href={href && (typeof href === "string" ? href : href.toString())}
      {...blockProps}
      ref={ref}
    >
      {iconID ? (
        <>
          <Icon icon={iconID} />
          <span>{children}</span>
        </>
      ) : (
        children
      )}
    </a>
  );
}
