import { forwardRef, type LegacyRef } from "react";

import styles from "./html.module.scss";

import { blockComponent, type IBlockProps } from "#components/meta";
import { Icon, type IIconID } from "#components/icons";

export interface IAnchourHTMLProps extends Omit<IBlockProps<"a">, "href"> {
  /**
   * Icon id to use in a link.
   */
  iconID?: IIconID;
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
  /**
   * Normalized `href` prop.
   */
  const finalHref = href && (typeof href === "string" ? href : href.toString());
  // show the href values
  const content = children ? children : finalHref;

  return (
    <a href={finalHref} {...blockProps} ref={ref}>
      {iconID ? (
        <>
          <Icon icon={iconID} />
          <span>{content}</span>
        </>
      ) : (
        content
      )}
    </a>
  );
}
