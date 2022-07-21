import {
  createContext,
  type ReactNode,
  useContext,
  useEffect,
  useState,
} from "react";

import { type ILoggedInUserData } from "#codegen/schema/001_interfaces";
import { getLocalStoreItem } from "#store/local";
import { fetchAccountInfo } from "#api/account";

export interface IAccountInfo extends ILoggedInUserData {
  isLoggedIn: boolean;
}

const AccountContext = createContext<IAccountInfo | undefined>(undefined);

interface IAccountProviderProps {
  children: ReactNode;
}

export function AccountProvider({ children }: IAccountProviderProps) {
  const [accountInfo, changeAccountInfo] = useState<
    ILoggedInUserData | undefined
  >(undefined);

  useEffect(() => {
    (async () => {
      const isRegistered = getLocalStoreItem<boolean>("is_registered", false);
      if (!isRegistered) {
        return;
      }

      const accountData = getLocalStoreItem<ILoggedInUserData>("account");

      if (accountData) {
        changeAccountInfo(accountData);
        return;
      }

      let remoteAccountData = undefined;

      try {
        remoteAccountData = await fetchAccountInfo();
      } catch (error) {
        console.log(error);
        return;
      }

      changeAccountInfo(remoteAccountData);
    })();
  }, []);

  return (
    <AccountContext.Provider
      value={{ ...accountInfo, isLoggedIn: Boolean(accountInfo) }}
    >
      {{ children }}
    </AccountContext.Provider>
  );
}

export function useAccount() {
  const accountInfo = useContext(AccountContext);

  return accountInfo;
}
