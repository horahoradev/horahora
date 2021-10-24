import { useCallback, useEffect, useRef, useState } from "react";
import { Tag, Avatar, Button, Rate } from "antd";
import { Link, useParams } from "react-router-dom";
import dashjs from "dashjs";
import videojs from "video.js";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUserCircle } from "@fortawesome/free-solid-svg-icons";

import * as API from "./api";
import Header from "./Header";
import { UserRank } from "./api/types";
import VideoList from "./VideoList";

const VIDEO_WIDTH = 44;
const VIDEO_HEIGHT = (9 / 16) * VIDEO_WIDTH;

function VideoPlayer(props) {
  let { url } = props;
  let videoRef = useRef();
  useEffect(() => {
    let video = videoRef.current;
    if (video == null) return;
  }, [url, videoRef]);

  var url_without_mpd = url.slice(0, -4);
  return (
    <>
      <video
          id="my-player"
          className="bg-black w-full object-contain object-center z-0"
          style={{ height: `${VIDEO_HEIGHT}rem` }}
          controls
          preload="auto"
          data-setup='{}'>
        <source src={url_without_mpd} type="video/mp4"></source>

        <p className="vjs-no-js">
          To view this video please enable JavaScript, and consider upgrading to a
          web browser that
          <a href="https://videojs.com/html5-video-support/" target="_blank">
            supports HTML5 video
          </a>
        </p>
      </video>
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
  let { data, id, setRating } = props;

  function rate(rating) {
    if (id == 0) {
      // TODO: throw
      return;
    }
    API.postRating(id, rating);
    setRating(rating);
  }

  // FIXME: new API endpoint
  return (
    <div className="bg-white border max-h-screen">
      <VideoPlayer url={data.MPDLoc} />
      <div className="p-4">
        <div>
          <span className="text-lg font-bold">{data.Title}</span>
          <span className="float-right">
            <span className="font-bold">{data.Views}</span> Views
          </span>
          <div className="inline-block relative top-5 float-right left-16"><Rate allowHalf={true} value={data.Rating} onChange={rate}></Rate></div>
          <br />
          <span className="text-gray-600 text-xs">{data.UploadDate}</span>
        </div>
        <div className="my-4">
          <span className="text-xs font-bold mb-2">Tags</span>
          <div className="border px-2 py-1">
            {data.Tags && data.Tags.map((tag, idx) => {
              // TODO(ivan): add links to tags
              return (
                <div key={idx} className="my-1 inline-block">
                  <a href={`/?search=${tag}`}><Tag>{tag}</Tag></a>
                </div>
              );
            })}
          </div>
        </div>
        {data.L && data.L.Rank === UserRank.ADMIN && <VideoAdminControls data={data} />}
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
  const [rating, setRating] = useState(0.0);
  const [userData, setUserData] = useState(null);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let data = await API.getVideo(id);
      if (!ignore) setPageData(data);
      if (!ignore) setRating(data.Rating);

      let userData = await API.getUserdata();
      if (!ignore) setUserData(userData);
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [id, rating]);

  if (pageData == null) return null;

  return (
    <>
      <Header userData={userData} />
      <div className="flex justify-center mx-4">
        <div className="max-w-screen-lg w-screen my-6 z-0">
          <VideoView data={pageData} id={id} setRating={setRating}/>
        </div>
        <div className="inline-block w-44 align-top float-right">
          <VideoList videos={pageData.RecommendedVideos} />
        </div>
      </div>
    </>
  );
}

export default VideoPage;
