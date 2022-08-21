import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { getPostComments, getPost } from "#api/lib";
import { VideoView } from "#components/posts";
import { Page } from "#components/page";
import { IVideoDetail, PostArticle } from "#entities/post";

function VideosPage() {
  const router = useRouter();
  const [videoDetail, changeVideoDetail] = useState<IVideoDetail>();
  const [comments, setComments] = useState([]);
  const { query, isReady } = router;
  let video_id = Number(query.video_id);
  const title = !videoDetail ? "Video" : `Video "${videoDetail.Title}"`


  async function refreshComments() {
    let videoComments = await getPostComments(video_id);
    setComments(videoComments);
  }

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    if (!isReady) {
      return;
    }
    let ignore = false;

    let fetchData = async () => {
      let data = await getPost(video_id);
      if (!ignore) changeVideoDetail(data);

      await refreshComments();
    };

    fetchData();

    return () => {
      ignore = true;
    };
  }, [video_id, isReady]);

  if (videoDetail == null) return null;

  return (
    <Page title={title} heading={null}>
      <PostArticle video={videoDetail} headingLevel={2}/>
      <VideoView
        data={videoDetail}
        videoComments={comments}
        id={video_id}
        refreshComments={refreshComments}
      />
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
