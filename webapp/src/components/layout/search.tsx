import { useRouter } from "next/router";

import { FormClient, ISubmitEvent } from "#components/forms";
import { Search } from "#components/inputs";

export function GlobalSearch() {
  const router = useRouter();
  const formID = "global-search"

  async function handleSubmit(event: ISubmitEvent) {
    router.push()
  }

  return (
    <FormClient
      id={formID}
      onSubmit={handleSubmit}
    >
      <Search id={`${formID}-query`} name="query">Search</Search>
    </FormClient>
  );
}
