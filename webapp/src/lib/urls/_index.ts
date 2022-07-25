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
