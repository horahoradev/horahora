import { Pagination } from "antd";

// Pagesize is not used

function Paginatione(props) {
    const paginationData = props.paginationData;
    const onPageChange = props.onPageChange;

    function changePage(page, number_of_items) {
        onPageChange(page);
    }

    return (
      <Pagination current={paginationData.CurrentPage} onChange={changePage} pageSize={50} total={paginationData.NumberOfItems}/>
  )
}

export default Paginatione;
