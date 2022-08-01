import { useEffect, useState } from "react";
import { Avatar } from "antd";
import { UserOutlined } from "@ant-design/icons";
import { useRouter } from "next/router";

import { getUser, banUser, setUserMod, setUserAdmin } from "#api/index";
import { type IUserRank, UserRank } from "#api/types";
import Paginatione from "#components/pagination";
import { Page } from "#components/page";
import { Button } from "#components/buttons";
import { type IVideo } from "#codegen/schema/001_interfaces";
import { PostList } from "#components/entities/post";

// {"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":11,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}

interface IPageUserData {
  Username: string;
  banned: boolean;
  L: Record<string, unknown> & {
    rank: IUserRank;
  };
  Videos: IVideo[];
  UserID: number;
  PaginationData: Record<string, unknown>;
}

/**
 * @TODO split into several pages
 */
function UsersPage() {
  const router = useRouter();
  const { query, isReady } = router;
  const { profile_id } = query;
  // @ts-expect-error typing
  const [pageUserData, setPageUserData] = useState<IPageUserData>({});
  const [currPage, setPage] = useState(1);

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    if (!isReady) {
      return;
    }

    let ignore = false;

    let fetchData = async () => {
      // @ts-expect-error some types
      let pageUserData = await getUser(profile_id, currPage);
      if (!ignore) setPageUserData(pageUserData);
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [profile_id, currPage, isReady]);

  return (
    <Page title="Profile information">
      {/*lol oh no*/}
      {/* TODO: add user profile image*/}
      <div className={"flex justify-center mx-4"}>
        <Avatar shape="square" size={96} icon={<UserOutlined />} />
      </div>

      <p className={"flex justify-center text-black dark:text-white"}>
        {pageUserData.Username} <b> {pageUserData.banned && "(banned)"} </b>
      </p>
      {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && (
        <p className={"flex justify-center"}>
          <Button onClick={() => banUser(pageUserData.UserID)}>Ban</Button>
        </p>
      )}
      {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && (
        <p className={"flex justify-center"}>
          <Button onClick={() => setUserMod(pageUserData.UserID)}>
            Promote to mod
          </Button>
        </p>
      )}
      {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && (
        <p className={"flex justify-center"}>
          <Button onClick={() => setUserAdmin(pageUserData.UserID)}>
            Promote to admin
          </Button>
        </p>
      )}

      {/* TODO: no more copy pasta! */}
      <div>
        <PostList posts={pageUserData.Videos} />
        <Paginatione
          paginationData={pageUserData.PaginationData}
          onPageChange={setPage}
        />
      </div>
    </Page>
  );
}

export default UsersPage;
