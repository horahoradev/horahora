import { useCallback, useEffect, useRef, useState } from "react";
import { Table, Timeline, Button, Space} from "antd";


import * as API from "./api";
import Header from "./Header";



function ArchivalPage() {
    const [userData, setUserData] = useState(null);
    const [archivalSubscriptions, setArchivalSubscriptions] = useState([]);
    const [timelineEvents, setTimelineEvents] = useState([]);

    // I think this is a hack? looks okay to me though!
    const [timerVal, setTimerVal] = useState(0);

    function reloadPage() {
        setTimerVal(timerVal => timerVal + 1);
    }

    async function deleteArchivalRequest(download_id) {
        await API.deleteArchivalRequest(download_id);
        reloadPage();   
    }

    useEffect(() => {
        const interval = setInterval(() => {
         reloadPage(timerVal);
        }, 30000);

        return () => clearInterval(interval);
      }, []);

    function createNewArchival() {
        const url = document.getElementById('url').value;
        API.postArchival(url);
        let newList = [{'Url': url, 'ArchivedVideos': 0, 'CurrentTotalVideos': 0, 'BackoffFactor': 1}, ...archivalSubscriptions];
        setArchivalSubscriptions(newList);
    }

    // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
    useEffect(() => {
        let ignore = false;

        let fetchData = async () => {
            let userData = await API.getUserdata();
            if (!ignore) setUserData(userData);

            let subscriptionData = await API.getArchivalSubscriptions();
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
            timelineEvents.map((event, idx) => <Timeline.Item>{event.message}<br></br>{event.timestamp}</Timeline.Item>),
        ];
    }

    const columns = [
        {
            title: 'URL',
            dataIndex: 'Url',
            key: 'Url',
        },
        {
            title: 'Downloaded videos',
            'dataIndex': 'ArchivedVideos',
            key: 'ArchivedVideos',
        },
        {
            title: 'Total videos',
            'dataIndex': 'CurrentTotalVideos',
            key: 'CurrentTotalVideos',
        },
        {
            title: 'Last synced',
            'dataIndex': 'LastSynced',
            key: 'LastSynced',
        },
        {
            title: 'Days until next sync',
            'dataIndex': 'BackoffFactor',
            key: 'BackoffFactor',
        },
        {
            title: 'Actions',
            key: 'action',
            render: (text, record) => (
                <Space size="middle">
                  <Button className="background-blue" onClick={()=>deleteArchivalRequest(record.downloadID)}>Delete {record.downloadID}</Button>
                </Space>
              ),
        }
    ];

    const timelinTableCols = [
        {
            title: 'Timestamp',
            dataIndex: 'timestamp',
            key: 'timestamp',
        },
        {
            title: 'Event Message',
            'dataIndex': 'message',
            key: 'message',
        }
    ];

    return (
        <>
            <Header userData={userData} />
            <br></br>
            Type URL to archive here: <input type="text" id="url">
            </input>
            <p></p>
            <Button className="background-blue" onClick={createNewArchival}>Submit</Button>
            <p></p>

            <div className={"inline-block w-1/2 mr-40"}>
                <Table dataSource={archivalSubscriptions} tableLayout={"auto"} columns={columns}/>
            </div>

            <div className={"inline-block w-2/5	"}>
                <Table dataSource={timelineEvents} columns={timelinTableCols}/>
            </div>
            </>
    );
}

export default ArchivalPage;