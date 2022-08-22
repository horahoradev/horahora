import { useRouter } from "next/router";

import { FormClient, ISubmitEvent } from "#components/forms";
import { LinkInternal } from "#components/links";
import {
  IListUnorderedProps,
  ListItem,
  ListUnordered,
} from "#components/lists";
import { blockComponent } from "#components/meta";
import { IURLBuilder, type IPagination } from "#lib/pagination";

// eslint-disable-next-line
import styles from "./internal.module.scss";

export interface IPaginationInternalProps extends IListUnorderedProps {
  pagination: IPagination;
  urlBuilder: IURLBuilder;
}

export const PaginationInternal = blockComponent(styles.block, Component);

function Component({
  pagination,
  urlBuilder,
  ...blockProps
}: IPaginationInternalProps) {
  const { limit, totalCount } = pagination;
  let { currentPage, totalPages } = pagination;
  totalPages = totalPages ?? Math.ceil(totalCount / limit);
  currentPage = currentPage ?? totalPages;
  const isLastPage = currentPage === totalPages;
  const prevPage = currentPage - 1;
  const nextPage = currentPage + 1;

  return (
    <ListUnordered {...blockProps}>
      <ListItem>
        {currentPage === 1 ? (
          "..."
        ) : (
          <LinkInternal href={urlBuilder(1)}>1</LinkInternal>
        )}
      </ListItem>

      <ListItem>
        {prevPage <= 1 ? (
          "..."
        ) : (
          <LinkInternal href={urlBuilder(prevPage)}>{prevPage}</LinkInternal>
        )}
      </ListItem>

      <ListItem>{currentPage}</ListItem>

      <ListItem>
        {nextPage >= totalPages ? (
          "..."
        ) : (
          <LinkInternal href={urlBuilder(nextPage)}>{nextPage}</LinkInternal>
        )}
      </ListItem>

      <ListItem>
        {isLastPage ? (
          "..."
        ) : (
          <LinkInternal href={urlBuilder(totalPages)}>
            {totalPages}
          </LinkInternal>
        )}
      </ListItem>
    </ListUnordered>
  );
}

interface ICurrentPageProps {
  pagination: IPagination;
  urlBuilder: IURLBuilder;
}

function CurrentPage({ pagination, urlBuilder }: ICurrentPageProps) {
  const router = useRouter();
  async function handleSubmit(event: ISubmitEvent) {}

  return (
    <FormClient
      id={`pagination`}
      onSubmit={handleSubmit}
      isSubmitSection={false}
    ></FormClient>
  );
}
