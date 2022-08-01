import { Pagination } from "antd";

export interface IPaginationeProps {
  paginationData: {
    CurrentPage?: number;
    NumberOfItems?: number;
  };
  onPageChange: (page: number) => void;
}

/**
 * Pagesize is not used
 */
function Paginatione({ paginationData, onPageChange }: IPaginationeProps) {
  function changePage(page: number, number_of_items: number) {
    onPageChange(page);
  }

  // TODO: what do here? ask Ivan..........
  if (!paginationData) {
    return <></>;
  }

  return (
    <Pagination
      current={paginationData.CurrentPage}
      onChange={changePage}
      pageSize={50}
      total={paginationData.NumberOfItems}
    />
  );
}

export default Paginatione;
