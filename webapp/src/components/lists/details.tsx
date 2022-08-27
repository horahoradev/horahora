import { type ReactNode } from "react";
import clsx from "clsx";

import { blockComponent, type IBlockProps } from "#components/meta";

// eslint-disable-next-line
import styles from "./details.module.scss";

export interface IDLProps extends IBlockProps<"dl"> {}
export interface IDSProps extends IBlockProps<"div"> {
  dKey?: ReactNode;
  dValue?: ReactNode;
  isHorizontal?: boolean;
}
export interface IDTProps extends IBlockProps<"dt"> {}
export interface IDDProps extends IBlockProps<"dd"> {}

/**
 * An html equivalent of dictionaries.
 */
export const DL = blockComponent(styles.block, DLComponent);
/**
 * A section of details list.
 */
export const DS = blockComponent(styles.section, DSComponent);
/**
 * A key of the details list.
 */
export const DT = blockComponent(styles.term, DTComponent);
/**
 * The value of the details list.
 */
export const DD = blockComponent(styles.definition, DDComponent);

function DLComponent({ children, ...blockProps }: IDLProps) {
  return <dl {...blockProps}>{children}</dl>;
}

function DSComponent({
  dKey,
  dValue,
  isHorizontal = false,
  children,
  className,
  ...blockProps
}: IDSProps) {
  const blockClass = clsx(
    className,
    children ?? styles.section_keyValue,
    isHorizontal && styles.section_horizontal
  );

  return (
    <div className={blockClass} {...blockProps}>
      {children ?? (
        <>
          <DT>{dKey}:</DT>
          <DD>{dValue}</DD>
        </>
      )}
    </div>
  );
}

function DTComponent({ children, ...blockProps }: IDTProps) {
  return <dt {...blockProps}>{children ?? "Unknown"}</dt>;
}

function DDComponent({ children, ...blockProps }: IDDProps) {
  return <dd {...blockProps}>{children ?? "Unknown"}</dd>;
}
