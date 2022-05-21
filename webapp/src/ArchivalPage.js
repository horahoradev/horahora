import { useCallback, useEffect, useRef, useState } from "react";
import { Input, Tag, Table, Timeline, Button, Space} from "antd";
import { CheckOutlined, SyncOutlined } from '@ant-design/icons';


import * as API from "./api";
import Header from "./Header";
import Footer from "./Footer";



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

    async function retryArchivalRequest(download_id) {
        await API.retryArchivalRequest(download_id);
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
        let subs =  archivalSubscriptions ? archivalSubscriptions : [];
        let newList = [{'Url': url, 'ArchivedVideos': 0, 'CurrentTotalVideos': 0, 'BackoffFactor': 1}, ...subs];
        setArchivalSubscriptions(newList);
    }

    function Status(record) {
        if (record.ArchivedVideos == record.CurrentTotalVideos && record.CurrentTotalVideos != 0){
            return (
                <Tag color="green" className="p-1 text-base" icon={<CheckOutlined/>}>Complete</Tag>
            );
        } else if(record.CurrentTotalVideos == 0 || record.LastSynced == null) {
            return (
                <Tag color="blue" className="text-base" icon={<SyncOutlined spin/>}>Fetching...</Tag>
            );
        } else {
            return (
                <Tag color="blue" className="text-base" icon={<SyncOutlined spin/>}>  Archiving...</Tag>
            );
        }
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
            title: 'Status',
            key: 'Url',
            render: (text, record) => (
                <span>
                   {Status(record)}
                </span>
              ),
        },
        {
            title: 'URL',
            dataIndex: 'Url',
            key: 'Url',
        },
        {
            title: 'Last synced',
            'dataIndex': 'LastSynced',
            key: 'LastSynced',
        },
        // {
        //     title: 'Days until next sync',
        //     'dataIndex': 'BackoffFactor',
        //     key: 'BackoffFactor',
        // },
        {
            title: 'Downloaded',
            key: 'Downloaded',
            render: (text, record) => (
                <span>
                    <b>{record.ArchivedVideos + "/" + record.CurrentTotalVideos}</b> videos ({record.UndownloadableVideos} undownloadable)
                </span>
              ),
        },
        {
            title: 'Actions',
            key: 'action',
            render: (text, record) => (
                <Space size="middle">
                  <Button className="background-blue" onClick={()=>deleteArchivalRequest(record.DownloadID)}>Delete {record.DownloadID}</Button>
                  <Button className="background-blue" onClick={()=>retryArchivalRequest(record.DownloadID)}>Retry {record.DownloadID}</Button>
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
            <div className="mx-10 h-full">
                <div className="inline-block mr-5">
                    <b className="text-4xl m-0">Archives</b>
                    <h2 className="text-xl mb-5">View and manage your archives</h2>
                    <div>

                    <div>
                    <div className={"inline-block mr-5 w-2/5 align-bottom"}>
                        <div className="bg-gray-50 border-gray-50 border-b-4">
                            <Input type="text" className="w-4/5 font-black text-base" placeholder="Paste URL to archive here" id="url">
                            </Input>
                            <Button className="w-1/5 text-base" type="primary"  onClick={createNewArchival}>Submit</Button>
                        </div>
                            <Table dataSource={archivalSubscriptions} scroll={{y: 700}} className="align-bottom w-full" ellipsis={true} columns={columns}/>
                    </div>
                    <div className="h-full inline-block w-2/5">

                        <Table dataSource={timelineEvents} className="align-bottom w-full" scroll={{y: 700}} ellipsis={true} columns={timelinTableCols}/>
                    </div>
                </div>
                </div>
                </div>

            </div>
                  </>
    );
}

export default ArchivalPage;