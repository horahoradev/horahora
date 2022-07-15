import { blockComponent, type IBlockProps } from "#components/meta";

export interface IFieldsetProps extends IBlockProps<"fieldset"> {}
export interface ILegendProps extends IBlockProps<"legend"> {}

export const Fieldset = blockComponent(undefined, FieldsetComponent);
export const Legend = blockComponent(undefined, LegendComponent);

function FieldsetComponent({ children, ...blockProps }: IFieldsetProps) {
  return <fieldset {...blockProps}>{children}</fieldset>;
}

function LegendComponent({ children, ...blockProps }: ILegendProps) {
  return <legend {...blockProps}>{children}</legend>;
}
