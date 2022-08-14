import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { fetchHome, type IHomeData } from "#api/lib";
import Paginatione from "#components/pagination";
import { Page } from "#components/page";
import { PostList } from "#components/entities/post";

export function HomePage() {
  const router = useRouter();
  const { query, isReady } = router;
  const [pageData, setPageData] = useState<IHomeData | null>(null);
  const [currPage, setPage] = useState(1);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    if (!isReady) {
      return;
    }

    let ignore = false;

    let fetchData = async () => {
      const { order, category, search } = query;

      try {
        let data: IHomeData = await fetchHome(
          currPage,
          search as string,
          order as string,
          category as string
        );
        if (!ignore) setPageData(data);
      } catch (error) {
        console.error(error);
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [currPage, isReady]);

  return (
    <Page title="Videos">
      <p>
        Number of videos: {pageData ? pageData.PaginationData.NumberOfItems : 0}
      </p>
      <PostList posts={pageData ? pageData.Videos : []} />
      <Paginatione
        paginationData={pageData ? pageData.PaginationData : {}}
        onPageChange={setPage}
      />
    </Page>
  );
}

export default HomePage;
