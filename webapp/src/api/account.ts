import { type IAccountClient } from "#lib/account";
import { apiFetch } from "#lib/fetch";

export async function fetchAccountInfo() {
  const accountInfo = await apiFetch<IAccountClient>(
    { pathname: "/currentuserprofile/" },
    {
      method: "GET",
      baseErrorMessage: "Failed fetch account info",
    }
  );

  return accountInfo;
}
