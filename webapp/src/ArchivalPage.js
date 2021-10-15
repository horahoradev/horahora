import { useCallback, useEffect, useRef, useState } from "react";
import { Table, Timeline, Button } from "antd";


import * as API from "./api";
import Header from "./Header";


function ArchivalPage() {
    const [userData, setUserData] = useState(null);
    const [archivalSubscriptions, setArchivalSubscriptions] = useState([]);
    const [timelineEvents, setTimelineEvents] = useState([]);

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
    }, []);

    let timelineElements = [];
    if (timelineEvents) {
        timelineElements = [
            timelineEvents.map((event, idx) => <Timeline.Item>{event.message}<br></br>{event.timestamp}</Timeline.Item>),
        ];
    }

    const columns = [
        {
            title: 'URL',
            dataIndex: 'url',
            key: 'url',
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

            <div className={"inline-block float-right   "}>
                <Timeline mode={"alternate"}>
                    {timelineElements}
                </Timeline>
            </div>


            </>
    );
}

export default ArchivalPage;