import { useCallback, useEffect, useRef, useState } from "react";
import { Tag, Avatar, Button } from "antd";
import { Link, useParams } from "react-router-dom";
import dashjs from "dashjs";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUserCircle } from "@fortawesome/free-solid-svg-icons";

import * as API from "./api";
import Header from "./Header";
import { UserRank } from "./api/types";

const VIDEO_WIDTH = 44;
const VIDEO_HEIGHT = (9 / 16) * VIDEO_WIDTH;

function VideoPlayer(props) {
  let { url } = props;
  let videoRef = useRef();
  useEffect(() => {
    let video = videoRef.current;
    if (video == null) return;
    let player = dashjs.MediaPlayer().create();
    player.initialize(video, url, false);
    return () => {
      player.destroy();
    };
  }, [url, videoRef]);

  return (
    <>
      <video
        ref={videoRef}
        className="bg-black w-full h-80 object-contain object-center"
        style={{ height: `${VIDEO_HEIGHT}rem` }}
        controls
      ></video>
    </>
  );
}

function VideoAdminControls(props) {
  let { data } = props;
  let [approvedVideo, setApprovedVideo] = useState(false);
  let approvingVideo = useRef(false);
  let approveVideo = useCallback(() => {
    if (approvingVideo.current) return;
    let run = async () => {
      await API.approveVideo(data.VideoID);
      setApprovedVideo(true);
    };
    approvingVideo.current = true;
    run().finally(() => {
      approvingVideo.current = false;
    });
  }, [data, approvingVideo, setApprovedVideo]);

  return (
    <>
      <hr />
      <div className="my-4">
        <Button type="primary" disabled={approvedVideo} onClick={approveVideo}>
          {approvedVideo ? "Approved" : "Approve"}
        </Button>
      </div>
    </>
  );
}

function VideoView(props) {
  let { data } = props;

  return (
    <div className="bg-white border" style={{ width: `${VIDEO_WIDTH}rem` }}>
      <VideoPlayer url={data.MPDLoc} />
      <div className="p-4">
        <div>
          <span className="text-lg font-bold">{data.Title}</span>
          <span className="float-right">
            <span className="font-bold">{data.Views}</span> Views
          </span>
          <span className="float-right">{/* TODO(ivan): rating */}</span>
          <br />
          <span className="text-gray-600 text-xs">{data.UploadDate}</span>
        </div>
        <div className="my-4">
          <span className="text-xs font-bold mb-2">Tags</span>
          <div className="border px-2 py-1">
            {data.Tags.map((tag, idx) => {
              // TODO(ivan): add links to tags
              return (
                <div key={idx} className="my-1 inline-block">
                  <Tag>{tag}</Tag>
                </div>
              );
            })}
          </div>
        </div>
        {data.L.Rank === UserRank.ADMIN && <VideoAdminControls data={data} />}
        <hr />
        <div className="my-4">
          <div className="flex justify-start items-center">
            <Link to={`/users/${data.AuthorID}`}>
              <Avatar
                size="large"
                icon={<FontAwesomeIcon icon={faUserCircle} />}
              />
            </Link>
            <div className="ml-3 text-center">
              {" "}
              <Link to={`/users/${data.AuthorID}`}>{data.Username}</Link>
            </div>
          </div>
          {/* TODO(ivan): THIS IS VERY OBVIOUSLY A SECURITY ISSUE! Remove this once you get a proper video description format going */}
          <div className="mt-4">
            <span dangerouslySetInnerHTML={{ __html: data.VideoDescription }} />
          </div>
        </div>
      </div>
    </div>
  );
}

function VideoPage() {
  let { id } = useParams();

  const [pageData, setPageData] = useState(null);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let data = await API.getVideo(id);
      if (!ignore) setPageData(data);
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [id]);

  if (pageData == null) return null;

  return (
    <>
      <Header userData={pageData.L} />
      <div className="flex justify-center mx-4">
        <div className="max-w-screen-lg w-screen my-6">
          <VideoView data={pageData} />
        </div>
      </div>
    </>
  );
}

export default VideoPage;
