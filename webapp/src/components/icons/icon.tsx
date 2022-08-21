import { type IconDefinition } from "@fortawesome/fontawesome-svg-core";
import {
  FontAwesomeIcon,
  type FontAwesomeIconProps,
} from "@fortawesome/react-fontawesome";

import { blockComponent } from "#components/meta";

// eslint-disable-next-line
import styles from "./icon.module.scss";

/**
 * An ID of the icon to be used by other components.
 */
export interface IIconID extends IconDefinition {}

export interface IIconProps extends FontAwesomeIconProps {}

export const Icon = blockComponent(styles.block, Component);

function Component({ ...blockProps }: IIconProps) {
  return <FontAwesomeIcon {...blockProps} />;
}
