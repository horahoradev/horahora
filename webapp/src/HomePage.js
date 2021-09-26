import { useEffect, useState } from "react";

import * as API from "./api";
import Header from "./Header";
import VideoList from "./VideoList";
import Paginatione from "./Pagination";

function HomePage() {
  const [pageData, setPageData] = useState(null);
  const [userData, setUserData] = useState(null);
  const [currPage, setPage] = useState(1);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let userData = await API.getUserdata();
      let data = await API.getHome(currPage);
      if (!ignore) setPageData(data);
      if (!ignore) setUserData(userData);
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [currPage]);

  if (pageData == null) return null;

  return (
    <>
      <Header userData={userData} />
      <div className="flex justify-center mx-4 min-h-screen">
        <div className="max-w-screen-lg w-screen my-6">
          <VideoList videos={pageData.Videos} />
          <Paginatione paginationData={pageData.PaginationData} onPageChange={setPage}/>
        </div>
      </div>
    </>
  );
}

export default HomePage;
