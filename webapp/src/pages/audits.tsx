import { useEffect, useState } from "react";
import { Table, Timeline, Button, Space, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArchive,
  faBars,
  faSearch,
  faSignOutAlt,
  faUser,
} from "@fortawesome/free-solid-svg-icons";
import type { ChangeEvent } from "react";

import { getAudits, getUserdata } from "#api/index";
import { Header } from "#components/header";
import { VideoList } from "#components/video-list";
import Paginatione from "#components/pagination";

function AuditsPage() {
  const [pageData, setPageData] = useState<{
    Events: Record<string, unknown>[];
    Length: number;
  } | null>(null);
  const [userData, setUserData] = useState(null);
  const [userID, setUserID] = useState(-1);
  const [currPage, setPage] = useState(1);

  const columns = [
    {
      title: "id",
      dataIndex: "ID",
      key: "ID",
    },
    {
      title: "User ID",
      dataIndex: "UserID",
      key: "UserID",
    },
    {
      title: "Message",
      dataIndex: "Message",
      key: "Message",
    },
    {
      title: "Timestamp",
      dataIndex: "Timestamp",
      key: "Timestamp",
    },
  ];

  function UserIDFFromSearch(e: ChangeEvent<HTMLInputElement>) {
    setUserID(Number(e.currentTarget.value));
  }

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let data = await getAudits(userID, currPage);
      if (!ignore) setPageData(data);

      let userData = await getUserdata();
      if (!ignore) setUserData(userData);
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [userID, currPage]);

  if (pageData == null) return null;

  return (
    <>
      <Header userData={userData} />

      <Input
        name="search"
        size="large"
        placeholder="Search for user ID"
        className="bg-white children:m-5 w-full text-black font-bold"
        onChange={(event) => {
          UserIDFFromSearch;
        }}
        prefix={
          <FontAwesomeIcon
            className="mr-1 text-gray-400 dark:text-white"
            icon={faSearch}
          />
        }
      />
      <Table
        className="bg-black"
        dataSource={pageData.Events}
        columns={columns}
        pagination={{
          current: currPage,
          onChange: setPage,
          pageSize: 50,
          total: pageData.Length,
        }}
      />
    </>
  );
}

export default AuditsPage;
