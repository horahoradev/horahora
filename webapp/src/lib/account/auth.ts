import { IAccountClient, IAccountInit } from "#codegen/schema/001_interfaces";

export async function registerAccount(accInit: IAccountInit): Promise<IAccountClient> {

}
export async function loginAccount(loginInfo: IAccLogin): Promise<IAccountClient> {}

export async function logoutAccount(): Promise<void> {}
