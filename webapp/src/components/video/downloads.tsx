import { useEffect, useState } from "react";
import { StompSubscription, type IMessage } from "@stomp/stompjs";
import { useMutex } from "react-context-mutex";

import { WSClient, WSConfig } from "#lib/fetch";
import { Page } from "#components/page";
import { fetchDownloadsInProgress } from "#api/archives";
import { CardList } from "#components/lists";
import { LoadingBar } from "#components/loading-bar";
import { DownloadCard, type IDownload } from "#entities/download";

export function ArchivalDownloadsPage() {
  const MutexRunner = useMutex();
  const mutex = new MutexRunner("messageHandler");
  const [conn, setConn] = useState<WSClient | null>(null);
  const [videoInProgressDataset, setVideoInProgressDataset] = useState<
    IDownload[] | null
  >([]);
  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);

  // TODO: currently connects every time the videos in progress changes
  useEffect(() => {
    var client = new WSClient(new WSConfig(), setConn);
    client.activate();

    return () => {
      client.deactivate();
    };
  }, []);

  // Get initial downloads in progress
  useEffect(() => {
    let unsub: StompSubscription[] = [];
    (async () => {
      let videos = await fetchDownloadsInProgress();
      for (var i = 0; i < (videos != null ? videos.length : 0); i++) {
        videos[i].progress = 0;
      }
      setVideoInProgressDataset(videos);

      for (var i = 0; i < (videos != null ? videos.length : 0); i++) {
        let videoID = videos[i].VideoID;
        if (conn != null) {
          let ret = conn.subscribe(`/topic/${videoID}`, processMessage, {
            "prefetch-count": String(100),
            ack: "client-individual",
            id: String(Math.random() * 1000),
          });
          unsub.push(ret);
        }
      }
    })();

    return () => {
      unsub.map((fn) => fn.unsubscribe());
    };
  }, [conn]);

  useEffect(() => {
    if (conn === null) {
      return;
    }
    // Videos in progress subscriptions
    conn.subscribe(
      `/topic/state`,
      function (message) {
        mutex.lock();
        let body = JSON.parse(message.body);
        if (body.Type == "deletion") {
          console.log(`Got delete ${body.Video.VideoID}`);
          setVideoInProgressDataset((videosInProg) => {
            if (videosInProg == null) {
              return videosInProg;
            }
            let dataset = JSON.parse(JSON.stringify(videosInProg));
            dataset = dataset.filter(
              // @ts-expect-error types
              (item) => item.VideoID != body.Video.VideoID
            );
            return dataset;
          });
          // Delete it from the list
        } else if (body.Type == "insertion") {
          console.log(`Got insert ${body.Video.VideoID}`);

          setVideoInProgressDataset((videosInProg) => {
            let dataset =
              videosInProg != null
                ? JSON.parse(JSON.stringify(videosInProg))
                : [];

            // // Does it already exist? If not, subscribe
            let videosID = dataset.filter(
              // @ts-expect-error types
              (item) => item.VideoID == body.Video.VideoID
            );
            if (videosID.length == 1) {
              // unsubscribing isn't important here
              conn.subscribe(`/topic/${body.Video.VideoID}`, processMessage, {
                "prefetch-count": String(100),
                ack: "client-individual",
              });
            }

            // Needed for upsert, filter it out if it's in there with a different dlStatus
            dataset = dataset.filter(
              // @ts-expect-error types
              (item) =>
                item.VideoID != body.Video.VideoID ||
                body.Video.DlStatus == item.DlStatus
            );
            body.Video.progress = 0;
            // If it's downloading, it goes at the beginning
            if (body.Video.DlStatus == "Downloading") {
              dataset.unshift(body.Video);
            } else {
              // else it goes at the end
              dataset.push(body.Video);
            }

            return dataset;
          });
        }
        message.ack();

        mutex.unlock();
      },
      {
        "prefetch-count": String(1),
        ack: "client-individual",
        id: String(Math.random() * 1000),
      }
    );
  }, [conn]);

  useEffect(() => {
    const interval = setInterval(() => {
      reloadPage();
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  function processMessage(message: IMessage) {
    mutex.lock();
    setVideoInProgressDataset((videosInProg: IDownload[] | null) => {
      if (videosInProg == null) {
        return videosInProg;
      }
      let dataset: IDownload[] = JSON.parse(JSON.stringify(videosInProg));
      let body = JSON.parse(message.body);
      let total_bytes = body.total_bytes || body.total_bytes_estimate;
      let progress =
        (100 * parseFloat(body.downloaded_bytes || total_bytes)) / total_bytes;
      let idx = dataset.findIndex(
        (video) => video.VideoID == body.info_dict.id
      );
      if (idx == -1) {
        return dataset;
      }
      dataset[idx].progress = progress;
      return dataset;
    });
    message.ack();
    mutex.unlock();
  }

  return (
    <Page title="Downloads in progress">
      <CardList>
        {!videoInProgressDataset ? (
          <LoadingBar />
        ) : (
          videoInProgressDataset.map((download, index) => (
            <DownloadCard key={index} download={download} />
          ))
        )}
      </CardList>
    </Page>
  );
}
