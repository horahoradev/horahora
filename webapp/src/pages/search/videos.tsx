import { useRouter } from "next/router";

import { FormClient } from "#components/forms";
import { Page } from "#components/page";

function VideoSearcPage() {
  const router = useRouter();
  const { isReady, query } = router;
  const { category, order, search } = query;

  return (
    <Page>
      <FormClient id="video-search" onSubmit={async () => {
        const params = new URLSearchParams([
          ["category", category],
          ["order", order],
          ["search", search],
        ]);
      }}></FormClient>
    </Page>
  );
}

export default VideoSearcPage;
