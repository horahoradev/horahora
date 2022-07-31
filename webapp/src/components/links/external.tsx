import { AnchourHTML, type IAnchourHTMLProps } from "./html";

import { blockComponent } from "#components/meta";

export interface ILinkExternalProps extends IAnchourHTMLProps {}

export const LinkExternal = blockComponent(undefined, Component);

function Component({
  referrerPolicy = "no-referrer",
  rel = "external noreferrer noopener",
  target = "_blank",
  children,
  ...blockProps
}: ILinkExternalProps) {
  return (
    <AnchourHTML
      referrerPolicy={referrerPolicy}
      rel={rel}
      target={target}
      {...blockProps}
    >
      {children}
    </AnchourHTML>
  );
}
