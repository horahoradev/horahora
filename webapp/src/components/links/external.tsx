import { AnchourHTML, type IAnchourHTMLProps } from "./html";

import { blockComponent } from "#components/meta";

export interface ILinkExternalProps extends IAnchourHTMLProps {}

export const LinkExternal = blockComponent(undefined, Component);

function Component({ children, ...blockProps }: ILinkExternalProps) {
  return <AnchourHTML {...blockProps}>{children}</AnchourHTML>;
}
