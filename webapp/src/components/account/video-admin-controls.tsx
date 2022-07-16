import { useCallback, useRef, useState } from "react";
import { useRouter } from "next/router";
import { Button } from "antd";

import {
  deleteVideo as apiDeleteVideo,
  approveVideo as apiApproveVideo,
} from "#api/index";

export interface IVideoAdminControls extends Record<string, unknown> {
  data: {
    VideoID: number;
  };
}

export function VideoAdminControls(props: IVideoAdminControls) {
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
