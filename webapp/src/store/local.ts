import { isError } from "#lib/errors";
import { toJSON, fromJSON } from "#lib/json";

/**
 * All local store values are JSON-encoded strings.
 */
export const LOCAL_STORAGE = {
  TEST: "__storage_test__",
  IS_REGISTERED: "is_registered",
  ACCOUNT: "account",
} as const;

export type ILocalStoreKey = typeof LOCAL_STORAGE[keyof typeof LOCAL_STORAGE];

export function getLocalStoreItem<Type = unknown>(
  storeName: ILocalStoreKey,
  defaultValue?: Type
): Type | undefined {
  const storageItem = localStorage.getItem(storeName);

  if (storageItem === null) {
    if (defaultValue) {
      setLocalStoreItem<Type>(storeName, defaultValue);
    }

    return defaultValue;
  }

  const item = fromJSON<Type>(storageItem);

  if (!item) {
    if (defaultValue) {
      setLocalStoreItem<Type>(storeName, defaultValue);
    }

    return defaultValue;
  }

  return item;
}

export function setLocalStoreItem<Type = unknown>(
  storeName: ILocalStoreKey,
  value: Type
) {
  try {
    const jsonValue = toJSON<Type>(value);
    localStorage.setItem(storeName, jsonValue);
  } catch (error) {
    if (!isError(error)) {
      throw error;
    }

    throw new Error(`Failed to set item "${storeName}" in \`localStorage\``, {
      cause: error,
    });
  }
}

export function deleteLocaleStoreItem(storeName: ILocalStoreKey) {
  localStorage.removeItem(storeName);
}
