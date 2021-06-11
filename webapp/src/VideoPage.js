import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import dashjs from "dashjs";

import { getVideo } from "./api";
import Header from "./Header";

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
        controls
      ></video>
    </>
  );
}

function VideoView(props) {
  let { data } = props;

  return (
    <div className="bg-white border max-w-xl w-full">
      <VideoPlayer url={data.MPDLoc} />
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
      let data = await getVideo(id);
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
