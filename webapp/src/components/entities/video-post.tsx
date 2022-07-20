import { Rate } from "antd";

import { LinkInternal } from "#components/links";
import { type IVideo } from "#types/entities";

export interface IVideoPostProps {
  inline: boolean;
  video: IVideo;
}

export function VideoPost({ video, inline }: IVideoPostProps) {
  return (
    <div
      className={
        inline ? "h-24 w-80 relative inline-block m-1" : "px-2 h-44 w-44 m-1"
      }
    >
      {/* @TODO: not make an entire component a link */}
      <LinkInternal href={`/videos/${video.VideoID}`}>
        <>
          <div className="rounded relative inline-block w-44">
            <img
              className="block w-44 h-24 object-cover object-center"
              alt={video.Title}
              src={`${video.ThumbnailLoc}`}
              onError={(e) => {
                (e.target as HTMLImageElement).onerror = null;
                (
                  e.target as HTMLImageElement
                ).src = `${video.ThumbnailLoc.slice(0, -6)}.jpg`;
              }}
            />
            {!inline && (
              <Rate
                className={"relative -mt-8 z-30 float-right"}
                allowHalf={true}
                disabled={true}
                value={video.Rating}
              ></Rate>
            )}
          </div>
          {/* TODO(ivan): deal with text truncation (hoping to have a multi-line text truncation,
                        which can't be done purely in css) */}
          <div
            className={
              inline
                ? "inline-block align-top h-44 inline-flex ml-2 justify-between  h-24 flex-col"
                : ""
            }
          >
            <div className="text-xs font-bold text-blue-500  py-1 text-black">
              {video.Title}
            </div>
            <div className="text-xs text-black">Views: {video.Views}</div>
            {inline && (
              <Rate
                className={"z-30"}
                allowHalf={true}
                disabled={true}
                value={video.Rating}
              ></Rate>
            )}
          </div>
        </>
      </LinkInternal>
    </div>
  );
}
