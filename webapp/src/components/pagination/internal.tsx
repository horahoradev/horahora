import { blockComponent, IChildlessBlockProps } from "#components/meta";

export interface IPaginationInternalProps extends IChildlessBlockProps<"div"> {}

export const PaginationInternal = blockComponent(undefined, Component);

function Component({ ...blockProps }: IPaginationInternalProps) {
  return <div {...blockProps}></div>;
}
