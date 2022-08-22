import { useCallback, useRef, useState } from "react";
import { useRouter } from "next/router";

import { Button } from "#components/buttons";
import { deletePost, approvePost } from "#api/lib";
import { IVideoDetail } from "#entities/post";
import { List, ListItem } from "#components/lists";

// eslint-disable-next-line
import styles from "./video-admin-controls.module.scss";

export interface IVideoAdminControls {
  data: IVideoDetail;
}

export function VideoAdminControls(props: IVideoAdminControls) {
  const router = useRouter();
  let { data } = props;
  // @TODO approval status check
  let [approvedVideo, setApprovedVideo] = useState(false);
  let approvingVideo = useRef(false);
  let deletingVideo = useRef(false);

  let deleteVideo = () => {
    if (deletingVideo.current) return;

    deletingVideo.current = true;

    let run = async () => {
      await deletePost(data.VideoID);
      deletingVideo.current = false;
      router.push("/");
    };
    run();
    // TODO: error future handler
  };

  let approveVideo = useCallback(() => {
    if (approvingVideo.current) return;

    let run = async () => {
      await approvePost(data.VideoID);
      setApprovedVideo(true);
    };
    approvingVideo.current = true;

    run().finally(() => {
      approvingVideo.current = false;
    });
    // TODO: error future handler
  }, [data, approvingVideo, setApprovedVideo]);

  return (
    <List isHorizontal className={styles.block}>
      <ListItem>
        <Button onClick={deleteVideo}>Delete video</Button>
      </ListItem>
      <ListItem>
        <Button disabled={approvedVideo} onClick={approveVideo}>
          {approvedVideo ? "Approved" : "Approve"}
        </Button>
      </ListItem>
    </List>
  );
}
