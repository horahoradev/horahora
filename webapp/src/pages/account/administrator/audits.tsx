import { useEffect, useState } from "react";
import { Table } from "antd";
import { faSearch } from "@fortawesome/free-solid-svg-icons";

import { getAudits } from "#api/lib";
import { Page } from "#components/page";
import { Icon } from "#components/icons";
import { Text } from "#components/inputs";
import { FormClient, IFormElements } from "#components/forms";

function AuditsPage() {
  const [pageData, setPageData] = useState<{
    Events: Record<string, unknown>[];
    Length: number;
  } | null>(null);
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

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    let fetchData = async () => {
      let data = await getAudits(userID, currPage);
      if (!ignore) setPageData(data);
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [userID, currPage]);

  if (pageData == null) return null;

  return (
    <Page title="Audit logs">
      <FormClient
        id="search"
        onSubmit={async (event) => {
          const elements = event.currentTarget
            .elements as IFormElements<"search">;
          setUserID(Number(elements["search"].value));
        }}
      >
        <Text id="user-search" name="search" maxLength={19}>
          <Icon icon={faSearch} /> Search for user ID
        </Text>
      </FormClient>

      <Table
        dataSource={pageData.Events}
        columns={columns}
        pagination={{
          current: currPage,
          onChange: setPage,
          pageSize: 50,
          total: pageData.Length,
        }}
      />
    </Page>
  );
}

export default AuditsPage;
