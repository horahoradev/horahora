import { createContext, useContext } from "react";

import { type ILoggedInUserData } from "#codegen/schema/001_interfaces";

export const AccountContext = createContext<ILoggedInUserData | undefined>(
  undefined
);

export function useAccount() {
  const accountInfo = useContext(AccountContext);

  return accountInfo;
}
