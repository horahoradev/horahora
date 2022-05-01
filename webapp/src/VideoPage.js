import { useCallback, useEffect, useRef, useState } from "react";
import { Tag, Avatar, Button, Input, Rate, Comment, List } from "antd";
import {Link, useHistory, useParams} from "react-router-dom";
import dashjs from "dashjs";
import videojs from "video.js";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faThumbsUp, faThumbsDown, faUserCircle, faUser} from "@fortawesome/free-solid-svg-icons";
import { useFormik } from "formik";

import * as API from "./api";
import Header from "./Header";
import { UserRank } from "./api/types";
import VideoList from "./VideoList";

const VIDEO_WIDTH = 44;
const VIDEO_HEIGHT = (9 / 16) * VIDEO_WIDTH;

function VideoPlayer(props) {
  let { url, next_video} = props;
  let videoRef = useRef();
  var url_without_mpd = url.slice(0, -4);

  useEffect(() => {
    let video = videoRef.current;
    if (video == null) return;
    video.setAttribute('src', url_without_mpd );
    video.load();
    video.play();

  }, [videoRef, url]);

  function set_on_end() {
      var video = document.getElementById("my-player");
      video.on('ended', () => {
        next_video();
      });
  }

  return (
    <>
      <video
          ref={videoRef}
          id="my-player"
          className="bg-black w-full object-contain object-center z-0"
          style={{ height: `${VIDEO_HEIGHT}rem` }}
          onEnded={next_video}
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
  let deletingVideo = useRef(false);
  let history = useHistory();

  let deleteVideo = () => { 
    if (deletingVideo.current) return;
    deletingVideo.current = true;
    let run = async () => {
      await API.deleteVideo(data.VideoID);
      deletingVideo.current = false;
      history.push("/");
    };
    run();
    // TODO: error future handler
  };

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
    // TODO: error future handler
  }, [data, approvingVideo, setApprovedVideo]);

  return (
    <>
      <hr />
      <div className="my-4">
        <Button type="primary" disabled={approvedVideo} onClick={approveVideo}>
          {approvedVideo ? "Approved" : "Approve"}
        </Button>
        <Button type="primary" onClick={deleteVideo}>
          Delete video
        </Button>
      </div>
    </>
  );
}

function VideoView(props) {
  let { data, id, setRating, next_video, videoComments, refreshComments} = props;

  // video_id, content (content of comment), and parent (parent comment id if a reply)
  let formik = useFormik({
    initialValues: {
      content: "",
      video_id: id,
      parent: "",
    },
    onSubmit: async (values) => {
      await API.postComment(values);
      await refreshComments();
    },
  });

  async function upvoteComment(commentID, has_upvoted) {
    await API.upvoteComment(commentID, !has_upvoted);
    await refreshComments();
  };

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
      <VideoPlayer url={data.MPDLoc} next_video={next_video}/>
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
            {!data.Tags ? "None" : data.Tags.map((tag, idx) => {
              // TODO(ivan): add links to tags
              return (
                <div key={idx} className="my-1 inline-block">
                  <a href={`/?search=${tag}`}><Tag>{tag}</Tag></a>
                </div>
              );
            })}
          </div>
        </div>
        {data.L && data.L.rank === UserRank.ADMIN && <VideoAdminControls data={data} />}
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
      <List
    className="comment-list"
    header={`${videoComments.length} replies`}
    itemLayout="horizontal"
    dataSource={videoComments}
    renderItem={item => (
      <li>
        <Comment
          actions={[<span>{item.upvote_count}</span>, <FontAwesomeIcon onClick={() => upvoteComment(item.id, item.user_has_upvoted) } className={item.user_has_upvoted ? "mr-1 text-green-400" : "mr-1 text-gray-400"} icon={faThumbsUp} />]}
          author={item.fullname}
          avatar={item.profile_picture_url}
          content={item.content}
          datetime={item.created}
        />
      </li>
    )}
  />,
   <form onSubmit={formik.handleSubmit}>
        <Input.Group>
        <Input
            name="content"
            values={formik.values.content}
            onChange={formik.handleChange}
            size="large"
            placeholder="<your comment here>"
            prefix={
              <FontAwesomeIcon className="mr-1 text-gray-400" icon={faUser} />
            }
          />
        </Input.Group>
        <Input.Group>
          <Button block type="primary" htmlType="submit" size="large">
            Submit
          </Button>
        </Input.Group>
    </form>
  </div>
  );
}

function VideoPage() {
  let history = useHistory();


  let { id } = useParams();

  const [pageData, setPageData] = useState(null);
  const [rating, setRating] = useState(0.0);
  const [comments, setComments] = useState([]);
  const [userData, setUserData] = useState(null);

  function navigate_to_next_video(){
    if (!pageData || !pageData.RecommendedVideos) return;
    history.push("/videos/" + pageData.RecommendedVideos[0].VideoID);
  };

  async function refreshComments(){
    let videoComments = await API.getComments(id);
    setComments(videoComments);
  }


  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let data = await API.getVideo(id);
      if (!ignore) setPageData(data);
      if (!ignore) setRating(data.Rating);

      let userData = await API.getUserdata();
      if (!ignore) setUserData(userData);

      await refreshComments();
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
          <VideoView data={pageData} videoComments={comments} id={id} refreshComments={refreshComments} setRating={setRating} next_video={navigate_to_next_video}/>
        </div>
        <div className="inline-block w-44 align-top float-right">
          <VideoList videos={pageData.RecommendedVideos} />
        </div>
      </div>
    </>
  );
}

export default VideoPage;
