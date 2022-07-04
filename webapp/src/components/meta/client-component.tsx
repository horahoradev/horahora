import { useEffect, useState } from "react";
import type { ReactElement } from "react";

import { IS_BROWSER } from "#environment/constants";

export interface IClientComponentProps {
  children: ReactElement;
}

/**
 * Renders the children only upon hydration.
 */
export function ClientComponent({ children }: IClientComponentProps) {
  const [isEnabled, enableComponent] = useState(false);

  useEffect(() => {
    IS_BROWSER && enableComponent(true);
  }, []);

  return <>{isEnabled ? children : <div>Loading...</div>}</>;
}
