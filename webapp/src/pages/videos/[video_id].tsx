import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { getComments, getVideo } from "#api/index";
import { VideoView } from "#components/posts";
import { type IVideoDetailed, type IVideo } from "#types/entities";
import { Page } from "#components/page";
import { VideoList } from "#components/video-list";

interface IPageData extends IVideoDetailed {
  RecommendedVideos: IVideo[];
}

function VideosPage() {
  const router = useRouter();
  const { query, isReady } = router;
  let video_id = Number(query.video_id);

  const [pageData, setPageData] = useState<IPageData | null>(null);
  const [rating, setRating] = useState(0.0);
  const [comments, setComments] = useState([]);

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

      await refreshComments();
    };

    fetchData();

    return () => {
      ignore = true;
    };
  }, [video_id, isReady]);

  if (pageData == null) return null;

  return (
    <Page title="Video">
      <VideoView
        data={pageData}
        videoComments={comments}
        id={video_id}
        refreshComments={refreshComments}
        setRating={setRating}
        rating={rating}
        next_video={navigate_to_next_video}
      />
      {/* <div className="ml-4 mt-2 w-100 align-top float-right">
        <VideoList
          videos={pageData.RecommendedVideos}
          title="Recommendations"
          inline={true}
        />
      </div> */}
    </Page>
  );
}

export default VideosPage;
