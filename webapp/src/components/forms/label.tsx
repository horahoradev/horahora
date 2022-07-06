import { IBlockProps } from "#components/meta";

export interface ILabelProps extends IBlockProps<"label"> {}

export function Label({ children, ...blockProps }: ILabelProps) {
  return (
    <label
      className="inline-block text-black dark:text-white py-1 cursor-pointer"
      {...blockProps}
    >
      {children}
    </label>
  );
}
