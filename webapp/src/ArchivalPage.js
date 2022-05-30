import { useCallback, useEffect, useRef, useState } from "react";
import { Input, Tag, Table, Timeline, Progress, Button, Space} from "antd";
import { CheckOutlined, SyncOutlined,  } from '@ant-design/icons';
import * as Stomp from '@stomp/stompjs';

import * as API from "./api";
import Header from "./Header";
import cloneDeep from 'lodash/cloneDeep';


let id =  Math.floor(Math.random() * 1000);
function ArchivalPage() {
    const [userData, setUserData] = useState(null);
    const [archivalSubscriptions, setArchivalSubscriptions] = useState([]);
    const [timelineEvents, setTimelineEvents] = useState([]);
    const [videoInProgressDataset, setVideoInProgressDataset] = useState([]);
    const [conn, setConn] = useState(null);
    const [initialDLDoneFlag, setDLDone] = useState(false);
    const latest = useRef(videoInProgressDataset);

    const [progress, setProgress] = useState(new Map());
    const [progFlag, setProgFlag] = useState(false);
    // TODO: currently connects every time the videos in progress changes
    useEffect(()=> {
            var client = new Stomp.Client({
                brokerURL: 'ws://localhost:15674/ws', // TODO
                connectHeaders: {
                  login: 'guest', // TODO
                  passcode: 'guest',
                },
                // debug: function (str) {
                //     console.log(str);
                //   },
                reconnectDelay: 5000,
                heartbeatIncoming: 4000,
                heartbeatOutgoing: 4000,
              });
              client.activate();

              client.onConnect = function(frame) {
                setConn(client);
              };
              
              client.onStompError = function (frame) {
                // Will be invoked in case of error encountered at Broker
                // Bad login/passcode typically will cause an error
                // Complaint brokers will set `message` header with a brief message. Body may contain details.
                // Compliant brokers will terminate the connection after any error
                console.log('Broker reported error: ' + frame.headers['message']);
                console.log('Additional details: ' + frame.body);
              };
              

            //   return () => client.deactivate();
    
    }, []);

      function processMessage(message) {
        let body = JSON.parse(message.body);
        let total_bytes = body.total_bytes || body.total_bytes_estimate;
        let prog = 100 * parseFloat(body.downloaded_bytes || total_bytes) / total_bytes;
       
        var poggers = progress;
        poggers[body.info_dict.id] =  prog;
        setProgress(poggers);
        setProgFlag(Math.random() * 100);
        message.ack();
    }
    

    useEffect(async() => {
        
        let videos = await API.getDownloadsInProgress();
        for (var i = 0; i < (videos != null ? videos.length : 0); i++) {
            videos[i].progress = 0;
        }
        setVideoInProgressDataset(videos);
        setDLDone(true);
    }, [])

    // Get initial downloads in progress
    useEffect( () => {
        if (!initialDLDoneFlag){
            return
        }
        let unsub = [];
        for (var i = 0; i < (videoInProgressDataset != null ? videoInProgressDataset.length : 0); i++) {
            let videoID = videoInProgressDataset[i].VideoID;
            if (conn != null) {
                let ret = conn.subscribe(`/topic/${videoID}`, (message)=>processMessage(message), {'prefetch-count': 1, 'ack': 'client'});
                unsub.push(ret);
            }
        }

        return ()=> unsub.map((fn)=>fn.unsubscribe());

    }, [conn, initialDLDoneFlag]);


    useEffect(() => {
            // Videos in progress subscriptions
            conn != null && conn.subscribe('/topic/videosinprogress', function(message) {
                let body = JSON.parse(message.body);
                let videos = JSON.parse(JSON.stringify(videoInProgressDataset));
                if (body.Type == "deletion") {
                    console.log("Got delete")
                    videos = videos.filter((item)=>item.VideoID != message.video.VideoID);
                    // Delete it from the list
                } else if (body.Type =="insertion") {
                    console.log("Got insert");
                    // Does it already exist? If not, subscribe
                    let videosID = videos.filter((item)=>item.VideoID != message.video.VideoID);
                    if (videosID.length == 0) { // unsubscribing isn't important here
                        conn.subscribe(`/topic/${message.video.VideoID}`, processMessage, {'prefetch-count': 1, 'ack': 'client'});
                    }

                    // Needed for upsert, filter it out if it's in there with a different dlStatus
                    videos = videos.filter((item)=>item.VideoID != message.video.VideoID || message.video.DlStatus != item.DlStatus);
                    message.video.progress = 0;
                    videos.push(message.video);
                }
                setVideoInProgressDataset(videos);
            });
            return conn;

      }, [conn]);

    

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

    useEffect(()=> {
        let wow = JSON.parse(JSON.stringify(videoInProgressDataset));
        for (var i = 0; i < (videoInProgressDataset != null ? videoInProgressDataset.length : 0); i++) {
            wow[i].progress = progress[wow[i].VideoID] || 0;
        }
        setVideoInProgressDataset(wow);
    }, [progress, progFlag]);

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
                  <Button className="background-blue" onClick={()=>retryArchivalRequest(record.DownloadID)}>Retry {record.DownloadID}</Button>
                  <Button className="background-blue" onClick={()=>deleteArchivalRequest(record.DownloadID)}>Delete {record.DownloadID}</Button>
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

    const videoDLsCols = [
        {
            title: 'Video ID',
            dataIndex: 'VideoID',
            key: 'videoID',
        },
        {
            title: 'Website',
            'dataIndex': 'Website',
            key: 'website',
        },
        {
            title: 'Progress',
            key: 'progress',
            render: (text, record) => (
                 <Progress percent={Math.floor(record.progress)} size="small" />
              ),
            
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
                        <h2 className="text-xl text-black">Archival Events</h2>
                        <Table dataSource={timelineEvents} className="align-bottom w-full" scroll={{y: 700}} ellipsis={true} columns={timelinTableCols}/>
                    </div>
                    <div className="h-full inline-block w-4/5">
                        <h2 className="text-xl text-black">Videos Currently Being Downloaded</h2>
                        <Table dataSource={videoInProgressDataset} className="align-bottom w-full" scroll={{y: 700}} ellipsis={true} columns={videoDLsCols}/>
                    </div>
                </div>
                </div>
                </div>

            </div>
                  </>
    );
}

export default ArchivalPage;