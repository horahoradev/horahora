import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import axios from "axios";

import { getHome, getUserdata } from "#api/index";
import { Header } from "#components/header";
import { VideoList } from "#components/video-list";
import Paginatione from "#components/pagination";

interface IPageData {
  PaginationData: Record<string, unknown>;
  Videos: Record<string, unknown>;
}

export function HomePage() {
  const router = useRouter();
  const { query, isReady } = router;
  const [pageData, setPageData] = useState<IPageData | null>(null);
  const [userData, setUserData] = useState(null);
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
        let data = await getHome(
          currPage,
          search as string,
          order as string,
          category as string
        );
        if (!ignore) setPageData(data);
      } catch (error) {
        // Bad redirect if not authenticated
        if (axios.isAxiosError(error) && error.response!.status === 403) {
          router.push("/login");
        }
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [currPage, isReady]);

  useEffect(() => {
    let ignore = false;

    (async () => {
      let userData = await getUserdata();
      if (!ignore) setUserData(userData);
    })();

    return () => {
      ignore = true;
    };
  }, []);

  return (
    <>
      <Header userData={userData} />
      <div className="flex justify-center mx-4 min-h-screen py-4">
        <div className="max-w-screen-lg w-screen">
          <h1 className="bold text-2xl text-black dark:text-white">
            Number of videos:{" "}
            {pageData ? pageData.PaginationData.NumberOfItems : 0}
          </h1>
          <VideoList
            // @ts-expect-error types
            videos={pageData ? pageData.Videos : []}
          />
          <Paginatione
            paginationData={pageData ? pageData.PaginationData : []}
            onPageChange={setPage}
          />
        </div>
      </div>
    </>
  );
}

export default HomePage;
