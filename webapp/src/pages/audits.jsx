import { useEffect, useState } from "react";
import { useLocation } from "react-router";
import { Table, Timeline, Button, Space, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArchive,
  faBars,
  faSearch,
  faSignOutAlt,
  faUser,
} from "@fortawesome/free-solid-svg-icons";

import * as API from "../api";
import { Header } from "../components/header";
import { VideoList } from "../components/video-list";
import Paginatione from "../Pagination";

function AuditsPage() {
  const [pageData, setPageData] = useState(null);
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

  function UserIDFFromSearch(e) {
    setUserID(e.currentTarget.value);
  }

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let data = await API.getAudits(userID, currPage);
      if (!ignore) setPageData(data);

      let userData = await API.getUserdata();
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
        onChange={UserIDFFromSearch}
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
