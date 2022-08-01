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
import { UserRank } from "#api/types";

interface IAccountContext {
  account?: IAccount;
  isInProgress: boolean;
  register: (...args: Parameters<typeof registerAccount>) => Promise<void>;
  login: (...args: Parameters<typeof loginAccount>) => Promise<void>;
  logout: (...args: Parameters<typeof logoutAccount>) => Promise<void>;
}

interface IUseAccount extends IAccountContext {
  isRegistered: boolean;
  isAdmin: boolean;
}

// this is a typescript ritual because default value
// has to have the same type as the context
// but these functions can't operate outside of the context component.
const defaultContext: IAccountContext = {
  isInProgress: false,
  register: async () => {},
  login: async () => {},
  logout: async () => {},
};

const AccountContext = createContext<IAccountContext>(defaultContext);

export function AccountProvider({ children }: { children: ReactNode }) {
  const [account, changeAccount] = useState<IAccount | undefined>(undefined);
  const [isInProgress, switchProgress] = useState(true);

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
      try {
        if (!isRegistered()) {
          return;
        }

        const accountData = getAccount();

        if (accountData) {
          changeAccount(accountData);
          return;
        }

        const remoteAccountData = await fetchAccountInfo();
        changeAccount(remoteAccountData);
      } catch (error) {
        console.log(error);
        return;
      } finally {
        switchProgress(false);
      }
    })();
  }, []);

  return (
    <AccountContext.Provider
      value={{ account, isInProgress, register, login, logout }}
    >
      {children}
    </AccountContext.Provider>
  );
}

// Using a hook so every component wouldn't need to import the context object and `useContext()`
// to get access to it
export function useAccount(): IUseAccount {
  const { account, ...accContext } = useContext(AccountContext);
  const isRegistered = Boolean(account);
  const isAdmin = account?.rank === UserRank.ADMIN;

  return { account, ...accContext, isRegistered, isAdmin };
}
