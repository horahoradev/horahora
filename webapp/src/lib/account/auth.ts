import {
  type IAccountClient,
  type IAccountInit,
  type IAccountLogin,
} from "./types";

import {
  registerAccount as fetchRegisterAccount,
  loginAccount as fetchLoginAccount,
  logoutAccount as fetchLogoutAccount,
} from "#api/authentication";
import { fetchAccountInfo } from "#api/account";
import {
  setLocalStoreItem,
  LOCAL_STORAGE,
  deleteLocaleStoreItem,
} from "#store/local";

export async function registerAccount(
  accInit: IAccountInit
): Promise<IAccountClient> {
  const formParams = new URLSearchParams();

  formParams.set("username", accInit.username);
  formParams.set("password", accInit.password);
  formParams.set("email", accInit.email);

  await fetchRegisterAccount(formParams);

  setLocalStoreItem<boolean>(LOCAL_STORAGE.IS_REGISTERED, true);

  const newAccount = await fetchAccountInfo();

  setLocalStoreItem<IAccountClient>(LOCAL_STORAGE.ACCOUNT, newAccount);

  return newAccount;
}

export async function loginAccount(
  loginInfo: IAccountLogin
): Promise<IAccountClient> {
  const formParams = new URLSearchParams();

  formParams.set("username", loginInfo.username);
  formParams.set("password", loginInfo.password);

  await fetchLoginAccount(formParams);

  setLocalStoreItem<boolean>(LOCAL_STORAGE.IS_REGISTERED, true);

  const account = await fetchAccountInfo();

  setLocalStoreItem<IAccountClient>(LOCAL_STORAGE.ACCOUNT, account);

  return account;
}

export async function logoutAccount(): Promise<void> {
  try {
    await fetchLogoutAccount();
  } finally {
    deleteLocaleStoreItem(LOCAL_STORAGE.IS_REGISTERED);
    deleteLocaleStoreItem(LOCAL_STORAGE.ACCOUNT);
  }
}
