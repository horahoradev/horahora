import Link, { type LinkProps } from "next/link";
import {
  type MouseEventHandler,
  type HTMLAttributeAnchorTarget,
  type ReactNode,
} from "react";


import styles from "./internal.module.scss";
import { AnchourHTML } from "./html";

import { blockComponent } from "#components/meta";
import { IIconID } from "#components/icons";

export interface ILinkInternalProps extends Omit<LinkProps, "passHref"> {
  target?: HTMLAttributeAnchorTarget;
  iconID?: IIconID;
  className?: string;
  children?: ReactNode;
  onClick?: MouseEventHandler<HTMLAnchorElement>;
}

export const LinkInternal = blockComponent(styles.block, Component);

export function Component({
  iconID,
  target = "_self",
  className,
  onClick,
  children,
  ...blockProps
}: ILinkInternalProps) {

  return (
    <Link {...blockProps} passHref>
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
