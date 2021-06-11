import { useEffect, useState } from "react";

import { getHome } from "./api";
import Header from './Header';
import VideoList from './VideoList';
import Pagination from './Pagination';

function HomePage() {
  const [pageData, setPageData] = useState(null);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let data = await getHome();
      if (!ignore) setPageData(data);
    };

    fetchData();
    return () => {
      ignore = true;
    }
  }, [])

  if (pageData == null) return null;

  return <>
    <Header userData={pageData.L} />
    <div className="flex justify-center mx-4">
      <div className="max-w-screen-lg w-screen my-6">
        <Pagination pagination={pageData.PaginationData} />
        <VideoList videos={pageData.Videos} />
        <Pagination pagination={pageData.PaginationData} />
      </div>
    </div>
  </>
}

export default HomePage;
