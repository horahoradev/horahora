import { blockComponent, type IChildlessBlockProps } from "#components/meta";
import { type IPagination } from "#lib/pagination";

// eslint-disable-next-line
import styles from "./info.module.scss";

export interface IPaginationInfoProps extends IChildlessBlockProps<"p"> {
  pagination: IPagination;
}

export const PaginationInfo = blockComponent(styles.block, Component);

function Component({ pagination, ...blockProps }: IPaginationInfoProps) {
  const { limit, totalCount } = pagination;
  let { currentPage, totalPages } = pagination;
  totalPages = totalPages ?? Math.ceil(totalCount / limit);
  currentPage = currentPage ?? totalPages;
  const isLastPage = currentPage === totalPages;
  const currentMin = (currentPage - 1) * limit + 1;
  const currentMax = isLastPage ? totalCount : currentPage * limit;

  return (
    <p {...blockProps}>
      Showing <span>{currentMin}</span>-<span>{currentMax}</span> out of total{" "}
      <span>{totalCount}</span> items.
    </p>
  );
}
