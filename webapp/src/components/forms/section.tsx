import type { IBlockProps } from "#components/meta";

export interface IFormSectionProps extends IBlockProps<"div"> {}

/**
 * Generic separator between form inputs.
 */
export function FormSection({ children, ...blockProps }: IFormSectionProps) {
  return (
    <div className="py-2" {...blockProps}>
      {children}
    </div>
  );
}
