import { useEffect, useRef, type SyntheticEvent } from "react";

import { blockComponent, type IBlockProps } from "#components/meta";
import { LinkExternal } from "#components/links";

// eslint-disable-next-line
import styles from "./player.module.scss";

export interface IVideoPlayerProps extends IBlockProps<"div"> {
  url: string;
  next_video: (event: SyntheticEvent<HTMLVideoElement, Event>) => void;
}

export const VideoPlayer = blockComponent(styles.block, Component);

export function Component({
  url,
  next_video,
  ...blockProps
}: IVideoPlayerProps) {
  let videoRef = useRef<HTMLVideoElement>(null);
  var url_without_mpd = url.slice(0, -4);

  useEffect(() => {
    let video = videoRef.current;
    if (video == null) return;
    video.setAttribute("src", url_without_mpd);
    video.load();
    // video.play();
  }, [videoRef, url]);

  return (
    <div {...blockProps}>
      <video
        ref={videoRef}
        id="my-player"
        className={styles.player}
        onEnded={next_video}
        controls
        preload="auto"
        data-setup="{}"
      >
        <source src={url_without_mpd} type="video/mp4"></source>

        <p className="vjs-no-js">
          To view this video please enable JavaScript, and consider upgrading to
          a web browser that
          <LinkExternal
            href="https://videojs.com/html5-video-support/"
            target="_blank"
            rel="noreferrer"
          >
            supports HTML5 video
          </LinkExternal>
        </p>
      </video>
    </div>
  );
}
