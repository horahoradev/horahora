import { type IconDefinition } from "@fortawesome/fontawesome-svg-core";
import {
  FontAwesomeIcon,
  type FontAwesomeIconProps,
} from "@fortawesome/react-fontawesome";

import { blockComponent, type IChildlessBlockProps } from "#components/meta";

// eslint-disable-next-line
import styles from "./icon.module.scss";

export interface IIconProps extends IChildlessBlockProps<"span"> {
  icon: FontAwesomeIconProps["icon"];
}

/**
 * An ID of the icon to be used by other components.
 */
export interface IIconID extends IconDefinition {}

export const Icon = blockComponent(styles.block, Component);

function Component({ icon, ...blockProps }: IIconProps) {
  return (
    <span {...blockProps}>
      <FontAwesomeIcon className={styles.icon} icon={icon} />
    </span>
  );
}
