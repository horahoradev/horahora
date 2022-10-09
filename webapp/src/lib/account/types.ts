export type {
  IAccount,
  IAccountClient,
  IAccountInit,
  IAccountLogin,
} from "#codegen/schema/001_interfaces";

export const UserRank = {
  REGULAR: 0,
  TRUSTED: 1,
  ADMIN: 2,
} as const;

export type IUserRank = typeof UserRank[keyof typeof UserRank];
