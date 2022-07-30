export interface IVideoDetailed {
  VideoID: number;
  MPDLoc: string;
  Title: string;
  Views: string;
  UploadDate: string;
  Tags: string[];
  AuthorID: number;
  Username: string;
  VideoDescription: string;
  L: {
    rank: number;
  };
}

export interface IComment {
  upvote_count: number;
  id: number;
  user_has_upvoted: boolean;
  content: string;
  created: string;
  fullname: string;
}
