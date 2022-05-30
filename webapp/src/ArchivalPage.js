import { useCallback, useEffect, useRef, useState } from "react";
import { Input, Tag, Table, Timeline, Progress, Button, Space} from "antd";
import { CheckOutlined, SyncOutlined,  } from '@ant-design/icons';
import * as Stomp from '@stomp/stompjs';
import { useMutex } from 'react-context-mutex';

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
    const latest = useRef(videoInProgressDataset);

    const MutexRunner = useMutex();
    const mutex = new MutexRunner('messageHandler');

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
      mutex.lock();
      setVideoInProgressDataset(videosInProg => {
          if (videosInProg == null) {
              return videosInProg;
          }
        let dataset = JSON.parse(JSON.stringify(videosInProg))
        let body = JSON.parse(message.body);
        let total_bytes = body.total_bytes || body.total_bytes_estimate;
        let progress = 100 * parseFloat(body.downloaded_bytes || total_bytes) / total_bytes;
        let idx = dataset.findIndex((video)=>video.VideoID == body.info_dict.id);
        if (idx == -1 ) {
            return dataset;
        }
        dataset[idx].progress = progress;
        return dataset;
      });

      message.ack();
      mutex.unlock();
  }


    // Get initial downloads in progress
    useEffect( async () => {

        let videos = await API.getDownloadsInProgress();
        for (var i = 0; i < (videos != null ? videos.length : 0); i++) {
            videos[i].progress = 0;
        }
        setVideoInProgressDataset(videos);

        let unsub = [];
        for (var i = 0; i < (videoInProgressDataset != null ? videoInProgressDataset.length : 0); i++) {
            let videoID = videoInProgressDataset[i].VideoID;
            if (conn != null) {
                let ret = conn.subscribe(`/topic/${videoID}`, processMessage, {'prefetch-count': 1, 'ack': 'client-individual'});
                unsub.push(ret);
            }
        }


        return ()=> unsub.map((fn)=>fn.unsubscribe());

    }, [conn]);


    useEffect(() => {
            // Videos in progress subscriptions
            conn != null && conn.subscribe('/topic/videosinprogress', function(message) {
                mutex.lock();
                let body = JSON.parse(message.body);
                message.ack();
                if (body.Type == "deletion") {
                    console.log(`Got delete ${body.Video.VideoID}`);
                    setVideoInProgressDataset(videosInProg => {
                        if (videosInProg == null) { return videosInProg; }
                        let dataset = JSON.parse(JSON.stringify(videosInProg))
                        dataset = dataset.filter((item)=>item.VideoID != body.Video.VideoID);
                        return dataset;
                    });
                    // Delete it from the list
                } else if (body.Type =="insertion") {
                    console.log(`Got insert ${body.Video.VideoID}`);

                    setVideoInProgressDataset(videosInProg => {
                        if (videosInProg == null) { return videosInProg; }
                        let dataset = JSON.parse(JSON.stringify(videosInProg))

                        // // Does it already exist? If not, subscribe
                        let videosID = dataset.filter((item)=>item.VideoID == body.Video.VideoID);
                        if (videosID.length == 1) { // unsubscribing isn't important here
                            conn.subscribe(`/topic/${body.Video.VideoID}`, processMessage, {'prefetch-count': 1, 'ack': 'client-individual'});
                        }

                        // Needed for upsert, filter it out if it's in there with a different dlStatus
                        dataset = dataset.filter((item)=>item.VideoID != body.Video.VideoID || body.Video.DlStatus != item.DlStatus);
                        body.Video.progress = 0;
                        // If it's downloading, it goes at the beginning
                        if (body.Video.DlStatus == 3) {
                            dataset.unshift(body.Video);
                        } else {
                           // else it goes at the end
                           dataset.push(body.Video);
                        }
                        
                        return dataset;
                    });
                }
                mutex.unlock();
            }, {'prefetch-count': 1, 'ack': 'client'});
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
            title: 'Download Status',
            key: 'DlStatus',
            render: (text, record) => (
                 <space size="middle">
                     { record.DlStatus == "Queued" || record.DlStatus == 4 ? "Queued" : "Downloading"}
                 </space>
             ),   
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