import {
  FontAwesomeIcon,
  type FontAwesomeIconProps,
} from "@fortawesome/react-fontawesome";

export interface IIconProps extends FontAwesomeIconProps {}

export function Icon({ ...blockProps }: IIconProps) {
  return (
    <FontAwesomeIcon
      className="max-h-5 text-black dark:text-white"
      {...blockProps}
    />
  );
}
