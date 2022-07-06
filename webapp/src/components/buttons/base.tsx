import { IBlockProps } from "#components/meta";

export interface IBaseButtonProps extends IBlockProps<"button"> {}

/**
 * Base button component for shared styles
 * between standard and submit buttons.
 */
export function BaseButton({ children, ...blockProps }: IBaseButtonProps) {
  return <button className="bg-white dark:bg-black p-2 rounded-md" {...blockProps}>{children}</button>;
}
