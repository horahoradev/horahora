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
    <Link
      className="flex items-center flex-row gap-2 text-base text-black dark:text-white dark:hover:text-black"
      href={href}
      {...blockProps}
    >
      {children ? children : href}
    </Link>
  );
}
