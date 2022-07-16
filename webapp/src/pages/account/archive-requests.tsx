import { useEffect, useRef, useState } from "react";
import { Tag, Table, Timeline, Progress, Button, Space } from "antd";
import { CheckOutlined, SyncOutlined } from "@ant-design/icons";
import {
  Client as StompClient,
  StompSubscription,
  type IMessage,
} from "@stomp/stompjs";
import { useMutex } from "react-context-mutex";

import {
  getDownloadsInProgress,
  deleteArchivalRequest as apiDeleteArchivalRequest,
  retryArchivalRequest as apiRetryArchivalRequest,
  postArchival,
  getUserdata,
  getArchivalSubscriptions,
} from "#api/index";
import { Header } from "#components/header";
import { NewVideoForm } from "#components/posts";

let id = Math.floor(Math.random() * 1000);

function ArchivalPage() {
  const [userData, setUserData] = useState(null);
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
  const latest = useRef(videoInProgressDataset);

  const MutexRunner = useMutex();
  const mutex = new MutexRunner("messageHandler");

  // TODO: currently connects every time the videos in progress changes
  useEffect(() => {
    var client = new StompClient({
      // brokerURL: "ws://localhost:61614/ws",
      webSocketFactory: function () {
        return new WebSocket("ws://localhost/ws");
      },
      connectHeaders: {
        login: "guest", // TODO
        passcode: "guest",
      },
      reconnectDelay: 5000,
      heartbeatIncoming: 4000,
      heartbeatOutgoing: 4000,
      connectionTimeout: 5000,
    });

    client.onConnect = function (frame) {
      setConn(client);
    };

    client.onDisconnect = function (frame) {
      setConn(null);
    };

    client.onWebSocketError = async (error) => {
      console.log(`onWebSocketError ${JSON.stringify(error)}`, "WS");
    };

    client.onStompError = function (frame) {
      // Will be invoked in case of error encountered at Broker
      // Bad login/passcode typically will cause an error
      // Complaint brokers will set `message` header with a brief message. Body may contain details.
      // Compliant brokers will terminate the connection after any error
      console.log("Broker reported error: " + frame.headers["message"]);
      console.log("Additional details: " + frame.body);
    };
    client.activate();

    return () => {
      (async () => {
        client.deactivate();
      })();
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

  async function deleteArchivalRequest(download_id: number) {
    await apiDeleteArchivalRequest(download_id);
    reloadPage();
  }

  async function retryArchivalRequest(download_id: number) {
    await apiRetryArchivalRequest(download_id);
    reloadPage();
  }

  useEffect(() => {
    const interval = setInterval(() => {
      reloadPage();
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  async function createNewArchival(url: string) {
    postArchival(url);
    let subs = archivalSubscriptions ? archivalSubscriptions : [];
    let newList = [
      { Url: url, ArchivedVideos: 0, CurrentTotalVideos: 0, BackoffFactor: 1 },
      ...subs,
    ];
    setArchivalSubscriptions(newList);
  }

  function Status(record: {
    ArchivedVideos: number;
    CurrentTotalVideos: number;
    LastSynced: null;
  }) {
    if (
      record.ArchivedVideos == record.CurrentTotalVideos &&
      record.CurrentTotalVideos != 0
    ) {
      return (
        <Tag color="green" className="p-1 text-base" icon={<CheckOutlined />}>
          Complete
        </Tag>
      );
    } else if (record.CurrentTotalVideos == 0 || record.LastSynced == null) {
      return (
        <Tag color="blue" className="text-base" icon={<SyncOutlined spin />}>
          Fetching...
        </Tag>
      );
    } else {
      return (
        <Tag color="blue" className="text-base" icon={<SyncOutlined spin />}>
          {" "}
          Archiving...
        </Tag>
      );
    }
  }

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

  const columns = [
    {
      title: "Status",
      key: "Status",
      render: (
        text: string,
        record: {
          ArchivedVideos: number;
          CurrentTotalVideos: number;
          LastSynced: null;
        }
      ) => <span>{Status(record)}</span>,
    },
    {
      title: "URL",
      dataIndex: "Url",
    },
    {
      title: "Last synced",
      dataIndex: "LastSynced",
    },
    // {
    //     title: 'Days until next sync',
    //     'dataIndex': 'BackoffFactor',
    //     key: 'BackoffFactor',
    // },
    {
      title: "Downloaded",
      key: "Downloaded",
      render: (
        text: string,
        record: {
          ArchivedVideos: number;
          CurrentTotalVideos: number;
          LastSynced: null;
          UndownloadableVideos: number;
        }
      ) => (
        <span>
          <b>{record.ArchivedVideos + "/" + record.CurrentTotalVideos}</b>{" "}
          videos ({record.UndownloadableVideos} undownloadable)
        </span>
      ),
    },
    {
      title: "Actions",
      key: "Actions",
      render: (
        text: string,
        record: {
          ArchivedVideos: number;
          CurrentTotalVideos: number;
          LastSynced: null;
          DownloadID: number;
        }
      ) => (
        <Space size="middle">
          <Button
            className="background-blue"
            onClick={() => retryArchivalRequest(record.DownloadID)}
          >
            Retry {record.DownloadID}
          </Button>
          <Button
            className="background-blue"
            onClick={() => deleteArchivalRequest(record.DownloadID)}
          >
            Delete {record.DownloadID}
          </Button>
        </Space>
      ),
    },
  ];

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
    <>
      <Header userData={userData} />
      <br></br>
      <div className=" mx-10 h-full">
        <div className="inline-block mr-5">
          <b className="text-black dark:text-white text-4xl m-0">Archives</b>
          <h2 className="text-black dark:text-white text-xl mb-5">
            View and manage your archives
          </h2>
          <div>
            <div>
              <div className="inline-block bg-white dark:bg-black mr-5 w-2/5 align-bottom">
                <NewVideoForm onNewURL={createNewArchival} />
                <Table
                  // @ts-expect-error types
                  dataSource={archivalSubscriptions}
                  scroll={{ y: 700 }}
                  className="align-bottom w-full "
                  ellipsis={true}
                  // @ts-expect-error types
                  columns={columns}
                />
              </div>
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
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

export default ArchivalPage;
