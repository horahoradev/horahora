export interface IPagination {
  currentPage?: number;
  totalCount: number;
  limit: number;
}

export interface ICollection<ItemType = never> {
  pagination: IPagination;
  items: ItemType[];
}
