import clsx from "clsx";

import { blockComponent, type IBlockProps } from "#components/meta";

// eslint-disable-next-line
import styles from "./headings.module.scss";

const headingLevels = [1, 2, 3, 4, 5, 6] as const;
export type IHeadingLevel = typeof headingLevels[number];

/**
 * Not creating separate interfaces per header because they all share the same interface.
 */
export interface IHeadingProps extends IBlockProps<"h1"> {
  level?: IHeadingLevel;
}

export const Heading = blockComponent(styles.block, Component);

function Component({
  level = 2,
  className,
  children,
  ...blockProps
}: IHeadingProps) {
  const finalClassName = clsx(styles[`level${level}`], className);

  switch (level) {
    case 1: {
      return (
        <h1 className={finalClassName} {...blockProps}>
          {children}
        </h1>
      );
    }

    case 2: {
      return (
        <h2 className={finalClassName} {...blockProps}>
          {children}
        </h2>
      );
    }

    case 3: {
      return (
        <h3 className={finalClassName} {...blockProps}>
          {children}
        </h3>
      );
    }

    case 4: {
      return (
        <h4 className={finalClassName} {...blockProps}>
          {children}
        </h4>
      );
    }

    case 5: {
      return (
        <h5 className={finalClassName} {...blockProps}>
          {children}
        </h5>
      );
    }

    case 6: {
      return (
        <h6 className={finalClassName} {...blockProps}>
          {children}
        </h6>
      );
    }

    default: {
      throw new Error(`Illegal heading level of "${level}"`);
    }
  }
}
