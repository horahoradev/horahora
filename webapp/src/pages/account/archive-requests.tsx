import { useEffect, useState } from "react";
import { Tag, Table, Timeline, Progress, Button, Space } from "antd";
import {
  Client as StompClient,
  StompSubscription,
  type IMessage,
} from "@stomp/stompjs";
import { useMutex } from "react-context-mutex";

import {
  getDownloadsInProgress,
  getUserdata,
  getArchivalSubscriptions,
} from "#api/index";
import { WSClient, WSConfig } from "#lib/fetch";
import { Page } from "#components/page";

function ArchivalPage() {
  const [userData, setUserData] = useState();
  const [archivalSubscriptions, setArchivalSubscriptions] = useState<
    {
      Url: string;
      ArchivedVideos: number;
      CurrentTotalVideos: number;
      BackoffFactor: number;
    }[]
  >([]);
  const [timelineEvents, setTimelineEvents] = useState<
    { message: string; timestamp: string }[]
  >([]);
  const [videoInProgressDataset, setVideoInProgressDataset] = useState<
    { VideoID: number; progress: number }[] | null
  >([]);
  const [conn, setConn] = useState<StompClient | null>(null);

  const MutexRunner = useMutex();
  const mutex = new MutexRunner("messageHandler");

  // TODO: currently connects every time the videos in progress changes
  useEffect(() => {
    var client = new WSClient(new WSConfig(), setConn);
    client.activate();

    return () => {
      client.deactivate();
    };
  }, []);

  function processMessage(message: IMessage) {
    mutex.lock();
    setVideoInProgressDataset(
      (videosInProg: { VideoID: number; progress: number }[] | null) => {
        if (videosInProg == null) {
          return videosInProg;
        }
        let dataset: { VideoID: number; progress: number }[] = JSON.parse(
          JSON.stringify(videosInProg)
        );
        let body = JSON.parse(message.body);
        let total_bytes = body.total_bytes || body.total_bytes_estimate;
        let progress =
          (100 * parseFloat(body.downloaded_bytes || total_bytes)) /
          total_bytes;
        let idx = dataset.findIndex(
          (video) => video.VideoID == body.info_dict.id
        );
        if (idx == -1) {
          return dataset;
        }
        dataset[idx].progress = progress;
        return dataset;
      }
    );
    message.ack();
    mutex.unlock();
  }

  // Get initial downloads in progress
  useEffect(() => {
    let unsub: StompSubscription[] = [];
    (async () => {
      let videos = await getDownloadsInProgress();
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

  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);

  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  useEffect(() => {
    const interval = setInterval(() => {
      reloadPage();
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let userData = await getUserdata();
      if (!ignore) setUserData(userData);

      let subscriptionData = await getArchivalSubscriptions();

      // videos.map((video, idx) => video.progress = videoInProgressDataset && videoInProgressDataset[idx] ? videoInProgressDataset[idx].progress : 0);

      // TODO: diff downloads in progress vs old downloads state, and unsubscribe!
      if (!ignore) {
        setArchivalSubscriptions(subscriptionData.ArchivalRequests);
        setTimelineEvents(subscriptionData.ArchivalEvents);
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [timerVal]);

  let timelineElements = [];

  if (timelineEvents) {
    timelineElements = [
      timelineEvents.map((event, idx) => (
        <Timeline.Item key={idx}>
          {event.message}
          <br></br>
          {event.timestamp}
        </Timeline.Item>
      )),
    ];
  }

  const timelinTableCols = [
    {
      title: "Timestamp",
      dataIndex: "timestamp",
    },
    {
      title: "Event Message",
      dataIndex: "message",
    },
  ];

  const videoDLsCols = [
    {
      title: "Video ID",
      dataIndex: "VideoID",
    },
    {
      title: "Website",
      dataIndex: "Website",
    },
    {
      title: "Download Status",
      dataIndex: "DlStatus",
    },
    {
      title: "Progress",
      key: "Progress",
      render: (text: string, record: { progress: number }) => (
        <Progress percent={Math.floor(record.progress)} size="small" />
      ),
    },
  ];

  return (
    <Page>

      <div className="h-full inline-block w-2/5">
        <h2 className="text-xl text-black dark:text-white">
          Archival Events
        </h2>
        <Table
          dataSource={timelineEvents}
          className="align-bottom w-full"
          scroll={{ y: 700 }}
          // @ts-expect-error types
          ellipsis={true}
          columns={timelinTableCols}
        />
      </div>
      <div className="h-full inline-block w-4/5">
        <h2 className="text-xl text-black dark:text-white">
          Videos Currently Being Downloaded
        </h2>
        <Table
          // @ts-expect-error types
          dataSource={videoInProgressDataset}
          className="align-bottom w-full"
          scroll={{ y: 700 }}
          ellipsis={true}
          columns={videoDLsCols}
        />
      </div>
    </Page>
  );
}

export default ArchivalPage;
