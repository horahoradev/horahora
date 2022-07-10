import { AnchourHTML, IAnchourHTMLProps } from "./html";
import styles from "./local.module.scss";

import { blockComponent } from "#components/meta";

export interface ILinkLocalProps extends IAnchourHTMLProps {}

export const LinkLocal = blockComponent(styles.block, Component);

function Component({ children, ...blockProps }: ILinkLocalProps) {
  return <AnchourHTML {...blockProps}>{children}</AnchourHTML>;
}
