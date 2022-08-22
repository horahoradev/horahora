import { blockComponent, type IChildlessBlockProps } from "#components/meta";
import { type IPagination } from "#lib/pagination";

// eslint-disable-next-line
import styles from "./info.module.scss";

export interface IPaginationInfoProps extends IChildlessBlockProps<"p"> {
  pagination: IPagination;
}

export const PaginationInfo = blockComponent(styles.block, Component);

function Component({ pagination, ...blockProps }: IPaginationInfoProps) {
  const { limit, totalCount, currentPage } = pagination;
  const totalPages = Math.ceil(totalCount / limit);
  const current = currentPage ?? totalPages;
  const isLastPage = currentPage === totalPages;
  const currentMin = (current - 1) * limit + 1;
  const currentMax = isLastPage ? totalCount : current * limit;

  return (
    <p {...blockProps}>
      Showing <span>{currentMin}</span>-<span>{currentMax}</span> out of total{" "}
      <span>{totalCount}</span> items.
    </p>
  );
}
