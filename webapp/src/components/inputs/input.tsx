import { IChildlessBlockProps } from "#components/meta";

export interface IInputProps extends IChildlessBlockProps<"input"> {}

export function Input({ ...blockProps }: IInputProps) {
  return <input className="bg-white dark:bg-black p-2 rounded-lg" {...blockProps} />;
}
