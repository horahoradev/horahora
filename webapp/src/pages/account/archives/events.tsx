import { Table, Timeline } from "antd";
import { useEffect, useState } from "react";

import { Page } from "#components/page";
import { getArchivalRequests } from "#api/archives";
import { type IArchivalEvent } from "#codegen/schema/001_interfaces";

function ArchivalEventsPage() {
  const [timelineEvents, setTimelineEvents] = useState<IArchivalEvent[]>([]);
  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);
  let timelineElements = [];
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

  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let subscriptionData = await getArchivalRequests();

      // videos.map((video, idx) => video.progress = videoInProgressDataset && videoInProgressDataset[idx] ? videoInProgressDataset[idx].progress : 0);

      // TODO: diff downloads in progress vs old downloads state, and unsubscribe!
      if (!ignore) {
        setTimelineEvents(subscriptionData.ArchivalEvents);
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

  return (
    <Page>
      <h1>Archival Events</h1>
      <Table
        dataSource={timelineEvents}
        className="align-bottom w-full"
        scroll={{ y: 700 }}
        // @ts-expect-error types
        ellipsis={true}
        columns={timelinTableCols}
      />
    </Page>
  );
}

export default ArchivalEventsPage;
