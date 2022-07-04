import { IBlockProps } from "#components/meta";

export interface ITextAreaProps extends IBlockProps<"textarea"> {}
export function TextArea({ children, ...blockProps }: ITextAreaProps) {
  return (
    <textarea className="bg-white dark:bg-black p-2 rounded-lg" {...blockProps}>
      {children}
    </textarea>
  );
}
