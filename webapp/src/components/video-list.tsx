import { PostCard } from "#components/entities/post";
import { type IVideo } from "#codegen/schema/001_interfaces";

const VIDEO_ELEMENT_WIDTH = "w-44";

interface IVideoListProps {
  title?: string;
  inline?: boolean;
  videos: IVideo[];
}

export function VideoList({ videos, title, inline }: IVideoListProps) {
  let elements: JSX.Element[] = [];

  if (videos) {
    elements = [
      // @ts-expect-error add spread
      videos.map((video, idx) => (
        <PostCard key={video.VideoID} headingLevel={2} post={video} />
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
