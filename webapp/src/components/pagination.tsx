import { Pagination } from "antd";

interface IPaginationeProps extends Record<string, unknown> {
  onPageChange: (page: number) => void;
}

/**
 * Pagesize is not used
 */
function Paginatione(props: IPaginationeProps) {
  const paginationData = props.paginationData;
  const onPageChange = props.onPageChange;

  function changePage(page: number, number_of_items: number) {
    onPageChange(page);
  }

  // TODO: what do here? ask Ivan..........
  if (!paginationData) {
    return <></>;
  }

  return (
    <Pagination
      // @ts-expect-error `paginationData` shape
      current={paginationData.CurrentPage}
      onChange={changePage}
      pageSize={50}
      // @ts-expect-error `paginationData` shape
      total={paginationData.NumberOfItems}
    />
  );
}

export default Paginatione;
