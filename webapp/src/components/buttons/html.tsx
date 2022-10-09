import { blockComponent, type IBlockProps } from "#components/meta";
import { Icon, type IIconID } from "#components/icons";

// eslint-disable-next-line
import styles from "./html.module.scss";

export interface IButtonHTMLProps extends IBlockProps<"button"> {
  iconID?: IIconID;
}

export const ButtonHTML = blockComponent(styles.block, Component);

/**
 * Base button component for shared styles
 * between standard and submit buttons.
 */
function Component({
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
