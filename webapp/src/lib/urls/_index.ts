import { PUBLIC_URL } from "#environment/derived";

export { normalizeQueryKey } from "./normalize-query-key";

export class HorahoraURL extends URL {
  constructor(url: string | URL) {
    super(url, PUBLIC_URL.origin);
    this.host = PUBLIC_URL.host;
  }
}

export class ProfileURL extends HorahoraURL {
  constructor(profileID: number) {
    super(`/profile/${profileID}`);
  }
}

export class VideoURL extends HorahoraURL {
  constructor(videoID: number) {
    super(`/videos/${videoID}`);
  }
}

const categories = ["upload_date", "rating", "views", "my_ratings"] as const;
const orders = ["asc", "desc"] as const;
type ICategory = typeof categories[number];
type IOrder = typeof orders[number];

export class SearchURL extends HorahoraURL {
  constructor(
    query: string,
    category: ICategory = "upload_date",
    order: IOrder = "desc"
  ) {
    super(`/search`);

    this.searchParams.set("search", query);
    this.searchParams.set("category", category);
    this.searchParams.set("order", order);
  }
}
