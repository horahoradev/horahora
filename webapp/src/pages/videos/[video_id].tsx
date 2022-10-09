import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { getPost } from "#api/lib";
import { Page } from "#components/page";
import { IVideoDetail, PostArticle } from "#entities/post";

function VideosPage() {
  const router = useRouter();
  const [videoDetail, changeVideoDetail] = useState<IVideoDetail>();
  const { query, isReady } = router;
  let video_id = Number(query.video_id);
  const title = !videoDetail
    ? "Video"
    : `Video "${videoDetail.Title}" (${videoDetail.VideoID}) by "${videoDetail.Username}" (${videoDetail.AuthorID})`;

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    if (!isReady) {
      return;
    }
    let ignore = false;

    let fetchData = async () => {
      let data = await getPost(video_id);
      if (!ignore) changeVideoDetail(data);
    };

    fetchData();

    return () => {
      ignore = true;
    };
  }, [isReady, video_id]);

  if (videoDetail == null) return null;

  return (
    <Page title={title} heading={null}>
      {/* lolz */}
      {<div className="whiteblock"></div>}
      <PostArticle video={videoDetail} headingLevel={1} />
      {/*
        <VideoList
          videos={pageData.RecommendedVideos}
          title="Recommendations"
          inline={true}
        />
      */}
    </Page>
  );
}

export default VideosPage;
