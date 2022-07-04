import Link, { type LinkProps } from "next/link";
import { type ReactNode } from "react";

export interface ILinkInternalProps extends LinkProps {
  children?: ReactNode;
}

export function LinkInternal({
  href,
  children,
  ...blockProps
}: ILinkInternalProps) {
  return (
    <Link href={href} {...blockProps}>
      {children ? children : href}
    </Link>
  );
}
