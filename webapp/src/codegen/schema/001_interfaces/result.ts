/*
  This module was created by the codegen, do not edit it manually.
*/
/**
 * The environment variables used by Horahora project.
 */
export interface IEnvironmentVariables {
  /**
   * https://docs.docker.com/compose/reference/envvars/#compose_project_name
   */
  COMPOSE_PROJECT_NAME?: string;
  /**
   * https://docs.docker.com/compose/reference/envvars/#compose_file
   */
  COMPOSE_FILE?: string;
  /**
   * https://docs.docker.com/compose/reference/envvars/#compose_docker_cli_build
   */
  COMPOSE_DOCKER_CLI_BUILD?: 0 | 1;
  /**
   * Enable buildkit for docker.
   */
  DOCKER_BUILDKIT?: 0 | 1;
  /**
   * The storage bucket name for b2 or s3.
   * If you want to use S3, you need to include your aws credentials and config in $HOME/.aws. The config and credentials will be mounted into the relevant services at runtime.
   * See https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html for more information.
   */
  BUCKET_NAME?: string;
  /**
   * For non-local applications, should be the public address of the server (if minio) or storage origin.
   * This will be the public URL of your Backblaze bucket WITH NO TRAILING SLASH. E.g. for me it's: https://f002.backblazeb2.com/file/otomads for backblaze, or https://horahora-dev-otomads.s3-us-west-1.amazonaws.com for s3.
   */
  ORIGIN_FQDN?: string;
  STORAGE_BACKEND?: "b2" | "s3";
  /**
   * The API ID for your Backblaze account if using backblaze, otherwise blank
   */
  STORAGE_API_ID?: string;
  /**
   * The API key for your Backblaze account, otherwise blank
   */
  STORAGE_API_KEY?: string;
  STORAGE_ENDPOINT?: string;
  VPN_PROVIDER?: string;
  VPN_IP?: string;
  WIREGUARD_PRIVATE_KEY?: string;
  VPN_CITIES?: string;
  SOCKS_ADDR?: string;
  GOLANG_PROTOBUF_REGISTRATION_CONFLICT?: string;
  GRPCPort?: number;
  JWT_KEYPAIR: string;
  [k: string]: unknown;
}

/**
 * A hohora account as shown on client.
 */
export interface IAccountClient {
  userID: number;
  username: string;
  profile_picture_url?: string;
  rank?: number;
  banned?: boolean;
}

/**
 * An account on hohora.
 */
export interface IAccount {
  userID: number;
  username?: string;
  profile_picture_url?: string;
  rank?: number;
  banned?: boolean;
}

/**
 * Initializer for horahora account,
 */
export interface IAccountInit {
  username: string;
}

export interface IArchivalEvent {
  video_url?: string;
  parent_url?: string;
  message?: string;
  timestamp?: string;
}

export interface IArchivalRequest {
  UserID?: number;
  Url?: string;
  ArchivedVideos?: number;
  CurrentTotalVideos?: number;
  LastSynced?: string;
  BackoffFactor?: number;
  DownloadID?: number;
  UndownloadableVideos?: number;
}

export interface IAuditData {
  Length?: number;
  Events?: IAuditEvent[];
}

export interface IAuditEvent {
  ID?: number;
  UserID?: number;
  Message?: string;
  Timestamp?: string;
}

export interface ICommentData {
  id: number;
  created: string;
  content: string;
  fullname: string;
  profile_picture_url: string;
  upvote_count: number;
  user_has_upvoted: boolean;
  parent?: number;
}

export interface IPaginationData {
  NumberOfItems?: number;
  CurrentPage?: number;
}

export interface IProfileData {
  PaginationData?: IPaginationData;
  UserID?: number;
  Username?: string;
  ProfilePictureURL?: string;
  Videos?: IVideo[];
  Banned?: boolean;
  L?: IAccount;
}

/**
 * Schema for the upload form
 */
export interface IFileUpload {
  title?: string;
  description?: string;
  tags?: string[];
  /**
   * Video file.
   */
  "file[0]": string;
  /**
   * Thumbnail file.
   */
  "file[1]": string;
}

export interface IVideoDetail {
  Title?: string;
  MPDLoc?: string;
  Views?: number;
  Rating?: number;
  VideoID?: number;
  AuthorID?: number;
  Username?: string;
  UserDescription?: string;
  VideoDescription?: string;
  UserSubscribers?: number;
  ProfilePicture?: string;
  /**
   * should be a datetime
   */
  UploadDate?: string;
  Tags?: string[];
  RecommendedVideos?: IVideo[];
  L?: IAccount;
}

export interface IVideoInProgress {
  Website?: string;
  VideoID?: string;
  DlStatus?: string;
}

export interface IVideoMetadata {
  /**
   * The location of the DASH manifest
   */
  videoLoc?: string;
  videoTitle?: string;
  rating?: number;
  /**
   * Do I need this? probably not
   */
  authorName?: string;
  views?: number;
  videoID?: number;
  uploadDate?: string;
  description?: string;
  authorID?: number;
  tags?: string[];
}

export interface IVideo {
  Title?: string;
  VideoID?: number;
  Views?: number;
  AuthorID?: number;
  AuthorName?: string;
  ThumbnailLoc?: string;
  Rating?: number;
}
