import {
  createContext,
  type ReactNode,
  useContext,
  useEffect,
  useState,
  useCallback,
} from "react";

import {
  getAccount,
  isRegistered,
  loginAccount,
  logoutAccount,
  registerAccount,
  type IAccount,
} from "#lib/account";
import { fetchAccountInfo } from "#api/account";

interface IAccountContext {
  account?: IAccount;
  register: (...args: Parameters<typeof registerAccount>) => void;
  login: (...args: Parameters<typeof loginAccount>) => void;
  logout: (...args: Parameters<typeof logoutAccount>) => void;
}

const defaultContext: IAccountContext = {
  register: () => {},
  login: () => {},
  logout: () => {},
};

const AccountContext = createContext<IAccountContext>(defaultContext);

export function AccountProvider({ children }: { children: ReactNode }) {
  const [account, changeAccount] = useState<IAccount | undefined>(undefined);

  // dunno if `useCallback()` is needed
  // but react can struggle with referential equality of functions
  // created outside of the rendering tree.
  const register = useCallback(
    async (...args: Parameters<IAccountContext["register"]>): Promise<void> => {
      const newAccount = await registerAccount(...args);
      changeAccount(newAccount);
    },
    []
  );

  const login = useCallback(
    async (...args: Parameters<IAccountContext["login"]>): Promise<void> => {
      const account = await loginAccount(...args);
      changeAccount(account);
    },
    []
  );

  const logout = useCallback(async (): Promise<void> => {
    await logoutAccount();
    changeAccount(undefined);
  }, []);

  // initialize the context
  useEffect(() => {
    (async () => {
      if (isRegistered()) {
        return;
      }

      const accountData = getAccount();

      if (accountData) {
        changeAccount(accountData);
        return;
      }

      try {
        const remoteAccountData = await fetchAccountInfo();
        changeAccount(remoteAccountData);
      } catch (error) {
        console.log(error);
        return;
      }
    })();
  }, []);

  return (
    <AccountContext.Provider value={{ account, register, login, logout }}>
      {{ children }}
    </AccountContext.Provider>
  );
}

// Using a hook so every component wouldn't need to import the context and `useContext()`
// to get access to it.
export function useAccount() {
  return useContext(AccountContext);
}
