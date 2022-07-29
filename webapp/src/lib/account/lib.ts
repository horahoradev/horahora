import { IAccountClient } from "#codegen/schema/001_interfaces";
import { getLocalStoreItem, LOCAL_STORAGE } from "#store/local";

/**
 * Gets the local copy of the account. Throws an error if the copy is not present.
 */
export function getAccount(): IAccountClient {
  if (!isRegistered()) {
    throw new Error("Account is not registered.")
  }

  const account = getLocalStoreItem<IAccountClient>(LOCAL_STORAGE.ACCOUNT);

  if (!account) {
    throw new Error("Account is registered but doesn't have a local copy.")
  }

  return account;
}

export function isRegistered(): boolean {
  const localResult = getLocalStoreItem<boolean>(LOCAL_STORAGE.IS_REGISTERED);

  if (!localResult) {
    return false;
  }

  return true;
}
