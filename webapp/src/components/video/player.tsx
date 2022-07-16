import { useEffect, useRef, type SyntheticEvent } from "react";

import { LinkExternal } from "#components/links";

export interface IVideoPlayerProps extends Record<string, unknown> {
  url: string;
  next_video: (event: SyntheticEvent<HTMLVideoElement, Event>) => void;
}

const VIDEO_WIDTH = 44;
const VIDEO_HEIGHT = (9 / 16) * VIDEO_WIDTH;

export function VideoPlayer(props: IVideoPlayerProps) {
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
          <LinkExternal
            href="https://videojs.com/html5-video-support/"
            target="_blank"
            rel="noreferrer"
          >
            supports HTML5 video
          </LinkExternal>
        </p>
      </video>
    </>
  );
}
