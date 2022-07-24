import { useEffect, useState } from "react";
import { Avatar, Button } from "antd";
import { UserOutlined } from "@ant-design/icons";
import { useRouter } from "next/router";

import {
  getUser,
  getUserdata,
  banUser,
  setUserMod,
  setUserAdmin,
} from "#api/index";
import { Header } from "#components/header";
import { type IUserRank, UserRank } from "#api/types";
import { VideoList } from "#components/video-list";
import Paginatione from "#components/pagination";

// {"PaginationData":{"PathsAndQueryStrings":["/users/1?page=1"],"Pages":[1],"CurrentPage":1},"UserID":1,"Username":"【旧】【旧】電ǂ鯨","ProfilePictureURL":"/static/images/placeholder1.jpg","Videos":[{"Title":"YOAKELAND","VideoID":1,"Views":11,"AuthorID":0,"AuthorName":"【旧】【旧】電ǂ鯨","ThumbnailLoc":"http://localhost:9000/otomads/7feaa38a-1e10-11ec-a6c3-0242ac1c0004.thumb","Rating":0}]}

interface IPageUserData extends Record<string, unknown> {
  Username: string;
  banned: boolean;
  L: Record<string, unknown> & {
    rank: IUserRank;
  };
  Videos: Record<string, unknown>[];
  UserID: number;
}

/**
 * @TODO split into several pages
 */
function UsersPage() {
  const router = useRouter();
  const { query, isReady } = router;
  const { user_id } = query;

  const [userData, setUserData] = useState();
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
      let pageUserData = await getUser(user_id, currPage);
      if (!ignore) setPageUserData(pageUserData);

      let userData = await getUserdata();
      if (!ignore) setUserData(userData);
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [user_id, currPage, isReady]);

  return (
    <>
      <Header userData={userData} />
      <div className="pt-10">
        {/*lol oh no*/}
        {/* TODO: add user profile image*/}
        <div className={"flex justify-center mx-4"}>
          <Avatar shape="square" size={96} icon={<UserOutlined />} />
        </div>

        <h1 className={"flex justify-center text-black dark:text-white"}>
          {pageUserData.Username} <b> {pageUserData.banned && "(banned)"} </b>
        </h1>
        {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && (
          <h1 className={"flex justify-center"}>
            <Button type="primary" onClick={() => banUser(pageUserData.UserID)}>
              Ban
            </Button>
          </h1>
        )}
        {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && (
          <h1 className={"flex justify-center"}>
            <Button
              type="primary"
              onClick={() => setUserMod(pageUserData.UserID)}
            >
              Promote to mod
            </Button>
          </h1>
        )}
        {pageUserData.L && pageUserData.L.rank === UserRank.ADMIN && (
          <h1 className={"flex justify-center"}>
            <Button
              type="primary"
              onClick={() => setUserAdmin(pageUserData.UserID)}
            >
              Promote to admin
            </Button>
          </h1>
        )}

        {/* TODO: no more copy pasta! */}
        <div className="flex justify-center mx-4 min-h-screen">
          <div className="max-w-screen-lg w-screen my-6">
            <VideoList
              // @ts-expect-error some type
              videos={pageUserData.Videos}
            />
            <Paginatione
              paginationData={pageUserData.PaginationData}
              onPageChange={setPage}
            />
          </div>
        </div>
      </div>
    </>
  );
}

export default UsersPage;
