import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { fetchHome, type IHomeData } from "#api/lib";
import { PaginationInfo, PaginationLocal } from "#components/pagination";
import { Page } from "#components/page";
import { PostList } from "#entities/post";

export function HomePage() {
  const router = useRouter();
  const { query, isReady } = router;
  const [pageData, setPageData] = useState<IHomeData>();
  const [currPage, setPage] = useState(1);
  const { isRegistered } = useAccount();

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    if (!isReady) {
      return;
    }

    let fetchData = async () => {
      const { order, category, search } = query;

      try {
        let data: IHomeData = await fetchHome(
          currPage,
          search as string,
          order as string,
          category as string
        );
        setPageData(data);
      } catch (error) {
        console.error(error);
      }
    };

    fetchData();
  }, [currPage, isReady, isRegistered]);

  return (
    <Page title="Videos">
      {!pageData ? (
        <p>Loading...</p>
      ) : (
        <>
          <PaginationInfo
            pagination={{
              totalCount: pageData.PaginationData.NumberOfItems!,
              currentPage: pageData.PaginationData.CurrentPage,
            }}
          />
          <PostList posts={pageData ? pageData.Videos : []} />
          <PaginationLocal
            pagination={{
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
