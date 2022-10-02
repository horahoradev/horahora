import { blockComponent, IChildlessBlockProps } from "#components/meta";

export interface ILoadingBarProps extends IChildlessBlockProps<"span"> {}

export const LoadingBar = blockComponent(undefined, Component)

function Component({...blockProps}:ILoadingBarProps) {
  return <span {...blockProps}>No data</span>
}
