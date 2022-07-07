import { useRouter } from "next/router";

import { FormClient, type ISubmitEvent } from "#components/forms";
import { Page } from "#components/page";
import { Search } from "#components/inputs";
import { normalizeQueryKey } from "#lib/urls";

function VideoSearcPage() {
  const router = useRouter();
  const { isReady, query } = router;
  const category = normalizeQueryKey(query.category, { defaultValue: "" });
  const order = normalizeQueryKey(query.order, { defaultValue: "" });
  const search = normalizeQueryKey(query.search, { defaultValue: "" });

  async function handleSearch(event: ISubmitEvent) {
    if (!isReady) {
      return;
    }

    const params = new URLSearchParams([
      ["category", category],
      ["order", order],
      ["search", search],
    ]);
  }

  return (
    <Page>
      <FormClient id="video-search" onSubmit={handleSearch}>
        <Search></Search>
      </FormClient>
    </Page>
  );
}

export default VideoSearcPage;
