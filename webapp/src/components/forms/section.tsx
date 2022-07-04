import type { IBlockProps } from "#components/meta";

export interface IFormSectionProps extends IBlockProps<"div"> {}

/**
 * Generic separator between form inputs.
 */
export function FormSection({ children, ...blockProps }: IFormSectionProps) {
  return <div {...blockProps}>{children}</div>;
}
