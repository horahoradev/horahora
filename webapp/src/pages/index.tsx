import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { fetchHome, type IHomeData } from "#api/lib";
import Paginatione, {
  PaginationInfo,
  PaginationLocal,
} from "#components/pagination";
import { Page } from "#components/page";
import { PostList } from "#components/entities/post";

export function HomePage() {
  const router = useRouter();
  const { query, isReady } = router;
  const [pageData, setPageData] = useState<IHomeData>();
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
      {!pageData ? (
        <p>Loading...</p>
      ) : (
        <>
          <PaginationInfo
            pagination={{
              // totalCount: 10000,
              totalCount: pageData.PaginationData.NumberOfItems!,
              currentPage: pageData.PaginationData.CurrentPage,
            }}
          />
          <PostList posts={pageData ? pageData.Videos : []} />
          <PaginationLocal
            pagination={{
              // totalCount: 10000,
              totalCount: pageData.PaginationData.NumberOfItems!,
              currentPage: pageData.PaginationData.CurrentPage,
            }}
            onPageChange={async (page) => {
              setPage(page);
            }}
          />
        </>
      )}
    </Page>
  );
}

export default HomePage;
