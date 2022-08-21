import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { fetchProfile } from "#api/lib";
import Paginatione from "#components/pagination";
import { Page } from "#components/page";
import { PostList } from "#entities/post";
import { IProfileData, ProfileArticle } from "#entities/profile";
import { LoadingBar } from "#components/loading-bar";

// {"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":11,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}

/**
 * @TODO split into several pages
 */
function UsersPage() {
  const router = useRouter();
  const { query, isReady } = router;
  const [pageUserData, setPageUserData] = useState<IProfileData>();
  const [currPage, setPage] = useState(1);
  const profile_id = Number(
    Array.isArray(query.profile_id) ? query.profile_id[0] : query.profile_id
  );

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    if (!isReady) {
      return;
    }

    let ignore = false;

    let fetchProfileData = async () => {
      let pageUserData = await fetchProfile(profile_id, currPage);
      if (!ignore) setPageUserData(pageUserData);
    };

    fetchProfileData();
    return () => {
      ignore = true;
    };
  }, [profile_id, currPage, isReady]);

  return (
    <Page title="Profile information">
      {!pageUserData ? (
        <LoadingBar />
      ) : (
        <>
          <ProfileArticle profile={pageUserData} />
          {/* @TODO make it a part of the article */}
          <div>
            <PostList posts={pageUserData.Videos} />
            <Paginatione
              paginationData={pageUserData.PaginationData}
              onPageChange={setPage}
            />
          </div>
        </>
      )}
    </Page>
  );
}

export default UsersPage;
