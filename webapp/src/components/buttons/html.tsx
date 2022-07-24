import styles from "./base.module.scss";

import { blockComponent, type IBlockProps } from "#components/meta";
import { Icon, type IIconID } from "#components/icons";

export interface IButtonHTMLProps extends IBlockProps<"button"> {
  iconID?: IIconID;
}

export const ButtonHTML = blockComponent(styles.block, Component);

/**
 * Base button component for shared styles
 * between standard and submit buttons.
 */
export function Component({
  iconID,
  children,
  ...blockProps
}: IButtonHTMLProps) {
  return (
    <button {...blockProps}>
      {iconID ? (
        <>
          <Icon icon={iconID} />
          <span>{children}</span>
        </>
      ) : (
        children
      )}
    </button>
  );
}
