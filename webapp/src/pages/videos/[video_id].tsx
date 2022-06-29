import { useCallback, useEffect, useRef, useState } from "react";
import { Tag, Avatar, Button, Input, Rate, Comment, List } from "antd";
import Link from "next/link";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faThumbsUp,
  faThumbsDown,
  faUserCircle,
  faUser,
} from "@fortawesome/free-solid-svg-icons";
import { useFormik } from "formik";
import { UserOutlined } from "@ant-design/icons";
import { useRouter } from "next/router";
import type { SyntheticEvent } from "react";

import {
  getComments,
  getVideo,
  getUserdata,
  deleteVideo as apiDeleteVideo,
  approveVideo as apiApproveVideo,
  postComment,
  upvoteComment as apiUpvoteComment,
  postRating,
} from "#api/index";
import { UserRank } from "#api/types";
import { VideoList } from "#components/video-list";
import { Header } from "#components/header";

const VIDEO_WIDTH = 44;
const VIDEO_HEIGHT = (9 / 16) * VIDEO_WIDTH;

interface IPageData extends Record<string, unknown> {
  RecommendedVideos: Record<string, unknown>[];
}

function VideosPage() {
  const router = useRouter();
  const { query, isReady } = router;
  let video_id = Number(query.video_id);

  const [pageData, setPageData] = useState<IPageData | null>(null);
  const [rating, setRating] = useState(0.0);
  const [comments, setComments] = useState([]);
  const [userData, setUserData] = useState(null);

  function navigate_to_next_video() {
    if (!pageData || !pageData.RecommendedVideos) return;
    router.push("/videos/" + pageData.RecommendedVideos[0].VideoID);
  }

  async function refreshComments() {
    let videoComments = await getComments(video_id);
    setComments(videoComments);
  }

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    if (!isReady) {
      return;
    }
    let ignore = false;

    let fetchData = async () => {
      let data = await getVideo(video_id);
      if (data) setRating(data.Rating);
      if (!ignore) setPageData(data);

      let userData = await getUserdata();
      if (!ignore) setUserData(userData);

      await refreshComments();
    };

    fetchData();

    return () => {
      ignore = true;
    };
  }, [video_id, isReady]);

  if (pageData == null) return null;

  return (
    <>
      <Header userData={userData} />
      <div className="flex justify-center mx-4">
        <div className="w-screen my-6 z-0 min-w-400">
          <VideoView
            // @ts-expect-error typing
            data={pageData}
            videoComments={comments}
            id={video_id}
            refreshComments={refreshComments}
            setRating={setRating}
            rating={rating}
            next_video={navigate_to_next_video}
          />
        </div>
        <div className="ml-4 mt-2 w-100 align-top float-right">
          <VideoList
            // @ts-expect-error types
            videos={pageData.RecommendedVideos}
            title="Recommendations"
            inline={true}
          />
        </div>
      </div>
    </>
  );
}

interface IVideoPlayerProps extends Record<string, unknown> {
  url: string;
  next_video: (event: SyntheticEvent<HTMLVideoElement, Event>) => void;
}

function VideoPlayer(props: IVideoPlayerProps) {
  let { url, next_video } = props;
  let videoRef = useRef<HTMLVideoElement>(null);
  var url_without_mpd = url.slice(0, -4);

  useEffect(() => {
    let video = videoRef.current;
    if (video == null) return;
    video.setAttribute("src", url_without_mpd);
    video.load();
    video.play();
  }, [videoRef, url]);

  return (
    <>
      <video
        ref={videoRef}
        id="my-player"
        className="bg-black dark:bg-black w-full max-w-screen-lg object-contain object-center z-0"
        style={{ height: `${VIDEO_HEIGHT}rem` }}
        onEnded={next_video}
        controls
        preload="auto"
        data-setup="{}"
      >
        <source src={url_without_mpd} type="video/mp4"></source>

        <p className="vjs-no-js">
          To view this video please enable JavaScript, and consider upgrading to
          a web browser that
          <a
            href="https://videojs.com/html5-video-support/"
            target="_blank"
            rel="noreferrer"
          >
            supports HTML5 video
          </a>
        </p>
      </video>
    </>
  );
}

interface IVideoAdminControls extends Record<string, unknown> {
  data: {
    VideoID: number;
  };
}

function VideoAdminControls(props: IVideoAdminControls) {
  const router = useRouter();
  let { data } = props;
  let [approvedVideo, setApprovedVideo] = useState(false);
  let approvingVideo = useRef(false);
  let deletingVideo = useRef(false);

  let deleteVideo = () => {
    if (deletingVideo.current) return;

    deletingVideo.current = true;

    let run = async () => {
      await apiDeleteVideo(data.VideoID);
      deletingVideo.current = false;
      router.push("/");
    };
    run();
    // TODO: error future handler
  };

  let approveVideo = useCallback(() => {
    if (approvingVideo.current) return;

    let run = async () => {
      await apiApproveVideo(data.VideoID);
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
        <b className="text-2xl text-black dark:text-white ml-4">
          Admin controls:
        </b>
        <Button
          type="primary"
          className="mx-4"
          disabled={approvedVideo}
          onClick={approveVideo}
        >
          {approvedVideo ? "Approved" : "Approve"}
        </Button>
        <Button type="primary" className="mx-4" onClick={deleteVideo}>
          Delete video
        </Button>
      </div>
    </>
  );
}

interface IVideoViewProps extends Record<string, unknown> {
  id: number;
  data: {
    MPDLoc: string;
    Title: string;
    Views: string;
    UploadDate: string;
    Tags: string[];
    AuthorID: number;
    Username: string;
    VideoDescription: string;
    L: {
      rank: number;
    };
  };
  videoComments: {
    upvote_count: number;
    id: number;
    user_has_upvoted: boolean;
    content: string;
    created: string;
    fullname: string;
  }[];

  rating: number;
  refreshComments: () => Promise<unknown>;
  setRating: (rating: number) => void;
  next_video: () => void;
}

function VideoView(props: IVideoViewProps) {
  let {
    rating,
    data,
    id,
    setRating,
    next_video,
    videoComments,
    refreshComments,
  } = props;

  // video_id, content (content of comment), and parent (parent comment id if a reply)
  let formik = useFormik({
    initialValues: {
      content: "",
      video_id: id,
      parent: "",
    },
    enableReinitialize: true,
    onSubmit: async (values) => {
      await postComment(values);
      await refreshComments();
    },
  });

  async function upvoteComment(commentID: number, has_upvoted: boolean) {
    await apiUpvoteComment(commentID, !has_upvoted);
    await refreshComments();
  }

  async function rate(rating: number) {
    if (id == 0) {
      // TODO: throw
      return;
    }
    await postRating(id, rating);
    setRating(rating);
  }

  // FIXME: new API endpoint
  return (
    <div className="bg-white dark:bg-black border">
      <VideoPlayer url={data.MPDLoc} next_video={next_video} />
      <div className="px-4 pt-4">
        <div>
          <span className="text-lg font-bold text-black dark:text-white">
            {data.Title}
          </span>
          <span className="float-right">
            <span className="text-black dark:text-white">
              {data.Views} Views
            </span>
          </span>
          <div className="inline-block relative top-5 float-right left-16 mr-2">
            <Rate allowHalf={true} value={rating} onChange={rate}></Rate>
          </div>
          <br />
          <span className="text-gray-600 text-xs">{data.UploadDate}</span>
        </div>

        <div className="my-2">
          <span className="text-xs font-bold text-black dark:text-white mb-2">
            Tags
          </span>
          <div className="border px-2 py-1">
            {!data.Tags
              ? "None"
              : data.Tags.map((tag, idx) => {
                  // TODO(ivan): add links to tags
                  return (
                    <div key={idx} className="my-1 inline-block">
                      <Link href={`/?search=${tag}`}>
                        <Tag color="blue">{tag}</Tag>
                      </Link>
                    </div>
                  );
                })}
          </div>
        </div>

        <hr className="box-border w-full"></hr>
        <div className="my-4">
          <div>
            <span className="h-20 w-20 inline-block">
              <Link href={`/users/${data.AuthorID}`}>
                <Avatar
                  size={80}
                  icon={<FontAwesomeIcon icon={faUserCircle} />}
                />
              </Link>
            </span>
            <span className="ml-2 pl-1 mt-2 inline-block align-top">
              <Link href={`/users/${data.AuthorID}`}>
                <b className="font-black text-blue-500 text-xl">
                  {data.Username}
                </b>
              </Link>
              <h1 className="text-black dark:text-white">0 subscribers</h1>
            </span>
          </div>
          <div className="ml-20 pl-3 text-black dark:text-white">
            {/* enjoy your XSS, bro */}
            <span dangerouslySetInnerHTML={{ __html: data.VideoDescription }} />
          </div>
        </div>
      </div>
      {data.L && data.L.rank === UserRank.ADMIN && (
        // @ts-expect-error types
        <VideoAdminControls data={data}></VideoAdminControls>
      )}
      <hr></hr>
      <List
        bordered={false}
        split={false}
        className="comment-list"
        header={
          <h2 className="ml-4 mb-0 text-xl text-black dark:text-white">
            Comments ({videoComments.length})
          </h2>
        }
        itemLayout="horizontal"
        dataSource={videoComments}
        renderItem={(item) => (
          <li>
            <Comment
              className="border-0 text-black dark:text-white shadow-none"
              actions={[
                <span key={1} className="text-black dark:text-white text-bold">
                  {item.upvote_count}
                </span>,
                <FontAwesomeIcon
                  key={2}
                  onClick={() => upvoteComment(item.id, item.user_has_upvoted)}
                  className={
                    item.user_has_upvoted
                      ? "mr-1 text-green-400"
                      : "mr-1 text-gray-400"
                  }
                  icon={faThumbsUp}
                />,
              ]}
              author={
                <b className="text-black dark:text-white">{item.fullname}</b>
              }
              avatar={
                <Avatar
                  className="ml-4"
                  shape="square"
                  size={50}
                  icon={<UserOutlined />}
                />
              }
              content={item.content}
              datetime={item.created}
            />
          </li>
        )}
      />
      ,
      <form onSubmit={formik.handleSubmit}>
        <Input.Group>
          <Input
            name="content"
            // @ts-expect-error types
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

export default VideosPage;
