import { IAccountInit } from "#codegen/schema/001_interfaces";

/**
 * Account info needed to login.
 */
export interface IAccLogin
  extends Pick<IAccountInit, "username" | "password"> {}
