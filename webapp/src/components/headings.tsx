import { blockComponent } from "#components/meta/block-component";
import { type IBlockProps } from "#components/meta/types";

const headingLevels = [1, 2, 3, 4, 5, 6] as const;
export type IHeadingLevel = typeof headingLevels[number];

/**
 * Not creating separate interfaces per header because they all share the same interface.
 */
export interface IHeadingProps extends IBlockProps<"h1"> {
  level: IHeadingLevel;
}

export const Heading = blockComponent(undefined, Component);

function Component({ level, children, ...blockProps }: IHeadingProps) {
  switch (level) {
    case 1: {
      return <h1 {...blockProps}>{children}</h1>;
    }
    case 2: {
      return <h2 {...blockProps}>{children}</h2>;
    }
    case 3: {
      return <h3 {...blockProps}>{children}</h3>;
    }
    case 4: {
      return <h4 {...blockProps}>{children}</h4>;
    }
    case 5: {
      return <h5 {...blockProps}>{children}</h5>;
    }
    case 6: {
      return <h6 {...blockProps}>{children}</h6>;
    }

    default: {
      throw new Error(`Illegal heading level of "${level}"`);
    }
  }
}
