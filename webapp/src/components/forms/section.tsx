import { blockComponent, IBlockProps } from "#components/meta";

export interface IFormSectionProps extends IBlockProps<"div"> {}

/**
 * Generic separator between form inputs.
 */
export const FormSection = blockComponent(undefined, Component);

export function Component({ children, ...blockProps }: IFormSectionProps) {
  return <div {...blockProps}>{children}</div>;
}
