import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import axios from "axios";

import { getHome } from "#api/index";
import Paginatione from "#components/pagination";
import { Page } from "#components/page";
import { PostList } from "#components/entities/post";
import { type IVideo } from "#codegen/schema/001_interfaces";

interface IPageData {
  PaginationData: Record<string, unknown>;
  Videos: IVideo[];
}

export function HomePage() {
  const router = useRouter();
  const { query, isReady } = router;
  const [pageData, setPageData] = useState<IPageData | null>(null);
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
        let data: IPageData = await getHome(
          currPage,
          search as string,
          order as string,
          category as string
        );
        if (!ignore) setPageData(data);
      } catch (error) {
        // Bad redirect if not authenticated
        if (axios.isAxiosError(error) && error.response!.status === 403) {
          router.push("/authentication/login");
        }
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
