import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { Page } from "#components/page";
import { getArchivalEvents, getArchivalRequests } from "#api/archives";
import { type IArchivalEvent } from "#codegen/schema/001_interfaces";
import { CardList } from "#components/lists";
import { LoadingBar } from "#components/loading-bar";
import { EventCard } from "#entities/event";

function ArchivalEventsPage() {
  const [events, changeEvents] = useState<IArchivalEvent[]>();
  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);

  const router = useRouter();
  const { query, isReady } = router;
  let download_id = Number(query.download_id);

  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    if (!isReady) {
      return
    }

    let fetchData = async () => {
      let subscriptionData = await getArchivalEvents(download_id);

      // videos.map((video, idx) => video.progress = videoInProgressDataset && videoInProgressDataset[idx] ? videoInProgressDataset[idx].progress : 0);

      // TODO: diff downloads in progress vs old downloads state, and unsubscribe!
      if (!ignore) {
        changeEvents(subscriptionData.ArchivalEvents);
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [timerVal, isReady]);

  useEffect(() => {
    const interval = setInterval(() => {
      reloadPage();
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  return (
    <Page title="Archival Events">
      <CardList>
        {!events ? (
          <LoadingBar />
        ) : (
          events.map((event, index) => <EventCard key={index} event={event} />)
        )}
      </CardList>
    </Page>
  );
}

export default ArchivalEventsPage;
