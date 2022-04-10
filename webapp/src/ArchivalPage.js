import { useCallback, useEffect, useRef, useState } from "react";
import { Table, Timeline, Button } from "antd";


import * as API from "./api";
import Header from "./Header";



function ArchivalPage() {
    const [userData, setUserData] = useState(null);
    const [archivalSubscriptions, setArchivalSubscriptions] = useState([]);
    const [timelineEvents, setTimelineEvents] = useState([]);

    // I think this is a hack? looks okay to me though!
    const [timerVal, setTimerVal] = useState(0);

    function reloadPage() {
        setTimerVal(timerVal + 1);
    }
    setInterval(reloadPage, 30000);

    function createNewArchival() {
        const url = document.getElementById('url').value;
        API.postArchival(url);
        let newList = [...archivalSubscriptions, url];
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

            <div className={"inline-block"}>
                <Table dataSource={archivalSubscriptions} columns={columns}/>
            </div>

            <div>
                <Table dataSource={timelineEvents} columns={timelinTableCols}/>
            </div>
            </>
    );
}

export default ArchivalPage;