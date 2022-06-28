import Link from "next/link";
import { Rate } from "antd";

const VIDEO_ELEMENT_WIDTH = "w-44";

export function VideoList(props) {
  const { videos, title, inline } = props;

  let elements = [];
  if (videos) {
    elements = [
      videos.map((video, idx) => (
        <Video inline={inline} key={idx} video={video} />
      )),
    ];
  }

  // add padding elements so the items in the last row on this flexbox grid
  // get aligned with the other rows
  for (let i = 0; i < 10; i++) {
    elements.push(
      <div
        key={`_padding_${i}`}
        className={inline ? "w-64" : VIDEO_ELEMENT_WIDTH}
      />
    );
  }

  return (
    <div
      className={
        inline
          ? "flex-col my-4 rounded border p-4 bg-white dark:bg-stone-800 flex flex-wrap min-h-screen"
          : "my-4 rounded border p-4 bg-white dark:bg-stone-800 flex flex-wrap min-h-screen"
      }
    >
      {title && (
        <h1 className="text-black dark:text-white ml-1 text-xl">{title}</h1>
      )}
      {elements}
    </div>
  );
}

function Video(props) {
  const { video, inline } = props;

  return (
    <div
      className={
        inline ? "h-24 w-80 relative inline-block m-1" : "px-2 h-44 w-44 m-1"
      }
    >
      <Link href={`/videos/${video.VideoID}`}>
        <div className="rounded relative inline-block w-44">
          <img
            className="block w-44 h-24 object-cover object-center"
            alt={video.Title}
            src={`${video.ThumbnailLoc}`}
            onError={(e) => {
              e.target.onerror = null;
              e.target.src = `${video.ThumbnailLoc.slice(0, -6)}.jpg`;
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
      </Link>
    </div>
  );
}
