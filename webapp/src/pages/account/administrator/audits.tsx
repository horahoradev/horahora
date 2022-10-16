import { useEffect, useState } from "react";
import { faSearch } from "@fortawesome/free-solid-svg-icons";

import { getAudits } from "#api/lib";
import { Page } from "#components/page";
import { Icon } from "#components/icons";
import { Text } from "#components/inputs";
import { FormClient, IFormElements } from "#components/forms";
import { CardList } from "#components/lists";
import { LoadingBar } from "#components/loading-bar";
import { AuditCard, IAuditData } from "#entities/audit";
import { PaginationInfo, PaginationLocal } from "#components/pagination";

function AuditsPage() {
  const [pageData, setPageData] = useState<IAuditData>();
  const [userID, setUserID] = useState(-1);
  const [currPage, setPage] = useState(1);

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
      {!pageData ? (
        <LoadingBar />
      ) : (
        <>
          <PaginationInfo
            pagination={{
              totalCount: pageData.Length!,
              currentPage: currPage,
            }}
          />
          <CardList>
            {pageData.Events != null && pageData.Events.map((audit) => (
              <AuditCard key={audit.ID} audit={audit} />
            ))}
          </CardList>
          <PaginationLocal
            pagination={{
              totalCount: pageData.Length!,
              currentPage: currPage,
            }}
            onPageChange={async (page) => {
              setPage(page);
            }}
          />
        </>
      )}
    </Page>
  );
}

export default AuditsPage;
