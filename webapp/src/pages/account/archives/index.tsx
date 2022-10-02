import { useEffect, useState } from "react";

import { Page } from "#components/page";
import { NewVideoForm } from "#components/posts";
import { IArchivalRequest } from "#codegen/schema/001_interfaces";
import { getArchivalRequests } from "#api/archives";
import { CardList } from "#components/lists";
import { ArchivalDownloadsPage } from "#components/video";

import { LoadingBar } from "#components/loading-bar";
import { RequestCard } from "#entities/request";

function NewArchivePage() {
  const [requests, changeRequests] = useState<IArchivalRequest[]>();
  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let subscriptionData = await getArchivalRequests();

      // videos.map((video, idx) => video.progress = videoInProgressDataset && videoInProgressDataset[idx] ? videoInProgressDataset[idx].progress : 0);

      // TODO: diff downloads in progress vs old downloads state, and unsubscribe!
      if (!ignore) {
        changeRequests(subscriptionData.ArchivalRequests);
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [timerVal]);

  useEffect(() => {
    const interval = setInterval(() => {
      reloadPage();
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  async function createNewArchival(url: string) {
    const subs = requests ? requests : [];
    reloadPage();
  }

  return (
    <Page title="View and manage your archives">
      <NewVideoForm onNewURL={createNewArchival} />
      <ArchivalDownloadsPage></ArchivalDownloadsPage>
      <CardList>
        {!requests ? (
          <LoadingBar />
        ) : (
          requests.map((request) => (
            <RequestCard
              key={request.DownloadID}
              request={request}
              onUpdate={reloadPage}
            />
          ))
        )}
      </CardList>
    </Page>
  );
}

export default NewArchivePage;
