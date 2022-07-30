import Link, { type LinkProps } from "next/link";
import {
  type MouseEventHandler,
  type HTMLAttributeAnchorTarget,
  type ReactNode,
} from "react";

import { AnchourHTML } from "./html";

import { HorahoraURL } from "#lib/urls";
import { blockComponent } from "#components/meta";
import { IIconID } from "#components/icons";

// eslint-disable-next-line
import styles from "./internal.module.scss";

export interface ILinkInternalProps extends Omit<LinkProps, "passHref"> {
  href: string | HorahoraURL;
  target?: HTMLAttributeAnchorTarget;
  iconID?: IIconID;
  className?: string;
  children?: ReactNode;
  onClick?: MouseEventHandler<HTMLAnchorElement>;
}

export const LinkInternal = blockComponent(styles.block, Component);

export function Component({
  href,
  iconID,
  target = "_self",
  className,
  onClick,
  children,
  ...blockProps
}: ILinkInternalProps) {
  return (
    <Link
      href={href instanceof HorahoraURL ? href.toString() : href}
      {...blockProps}
      passHref
    >
      <AnchourHTML
        className={className}
        target={target}
        iconID={iconID}
        onClick={onClick}
      >
        {children}
      </AnchourHTML>
    </Link>
  );
}
