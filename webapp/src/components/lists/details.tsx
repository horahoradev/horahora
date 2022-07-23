import { blockComponent, type IBlockProps } from "#components/meta";

export interface IDLProps extends IBlockProps<"dl"> {}
export interface IDSProps extends IBlockProps<"div"> {}
export interface IDTProps extends IBlockProps<"dt"> {}
export interface IDDProps extends IBlockProps<"dd"> {}

/**
 * An html equivalent of dictionaries.
 */
export const DL = blockComponent(undefined, DLComponent);
/**
 * A section of details list.
 */
export const DS = blockComponent(undefined, DSComponent);
/**
 * A key of the details list.
 */
export const DT = blockComponent(undefined, DTComponent);
/**
 * The value of the details list.
 */
export const DD = blockComponent(undefined, DDComponent);

function DLComponent({ children, ...blockProps }: IDLProps) {
  return <dl {...blockProps}>{children}</dl>;
}

function DSComponent({ children, ...blockProps }: IDSProps) {
  return <div {...blockProps}>{children}</div>;
}

function DTComponent({ children, ...blockProps }: IDTProps) {
  return <dt {...blockProps}>{children}</dt>;
}

function DDComponent({ children, ...blockProps }: IDDProps) {
  return <dd {...blockProps}>{children}</dd>;
}
