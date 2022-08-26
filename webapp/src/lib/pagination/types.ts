export interface IPagination {
  currentPage?: number;
  totalPages?: number;
  totalCount: number;
  limit: number;
}

export interface ICollection<ItemType = never> {
  pagination: IPagination;
  items: ItemType[];
}

/**
 * A function which accepts a page number and returns the URL string for this page.
 */
export type IURLBuilder = (page: number) => string | URL
