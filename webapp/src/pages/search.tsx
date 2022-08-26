import { useRouter } from "next/router";
import { useEffect, useState } from "react";

import { Page } from "#components/page";
import { fetchHome, type IHomeData } from "#api/lib";
import { PostList } from "#entities/post";
import { PaginationInfo, PaginationLocal } from "#components/pagination";
import { FormClient, IFormElements, ISubmitEvent } from "#components/forms";
import { RadioGroup, Select, Text } from "#components/inputs";
import { LoadingBar } from "#components/loading-bar";

function PostSearch() {
  const router = useRouter();
  const { query, isReady } = router;
  const [pageData, setPageData] = useState<IHomeData | null>(null);
  const [currPage, setPage] = useState(1);
  const formID = "post-search";
  const title = "Search Posts";

  useEffect(() => {
    if (!isReady) {
      return;
    }

    let ignore = false;

    let fetchData = async () => {
      const { order, category, search } = query;

      try {
        let data: IHomeData = await fetchHome(
          currPage,
          search as string,
          order as string,
          category as string
        );
        if (!ignore) setPageData(data);
      } catch (error) {
        console.error(error);
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [currPage, isReady]);

  async function handleSubmit(event: ISubmitEvent) {
    if (!isReady) {
      return;
    }

    const elements = event.currentTarget.elements as IFormElements<
      "search" | "order" | "category"
    >;

    const search = elements["search"].value;
    const category = elements["category"].value;
    const order = elements["order"].value;

    let data: IHomeData = await fetchHome(
      currPage,
      search as string,
      order as string,
      category as string
    );
    setPageData(data);
  }

  return (
    <Page title={title}>
      <FormClient id={formID} onSubmit={handleSubmit}>
        <Text id={`${formID}-search`} name="search">
          Query
        </Text>
        <Select
          id={`${formID}-category`}
          name="category"
          options={[
            { title: "Upload Date", value: "upload_date" },
            { title: "Rating", value: "rating" },
            { title: "Views", value: "views" },
            { title: "My ratings", value: "my_ratings" },
          ]}
        >
          Category
        </Select>
        <RadioGroup
          name="order"
          options={[
            { id: `${formID}-order-asc`, title: "Ascending", value: "asc" },
            { id: `${formID}-order-desc`, title: "Descending", value: "desc" },
          ]}
        >
          Order
        </RadioGroup>
      </FormClient>
      {!pageData ? (
        <LoadingBar />
      ) : (
        <>
          <PaginationInfo
            pagination={{
              totalCount: pageData.PaginationData.NumberOfItems!,
              currentPage: pageData.PaginationData.CurrentPage,
            }}
          />
          <PostList posts={pageData ? pageData.Videos : []} />
          <PaginationLocal
            pagination={{
              totalCount: pageData.PaginationData.NumberOfItems!,
              currentPage: pageData.PaginationData.CurrentPage,
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

export default PostSearch;
