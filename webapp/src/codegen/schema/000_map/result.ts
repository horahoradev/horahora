/*
  This module was created by the codegen, do not edit it manually.
*/
export const environmentVariablesSchema = {
  $id: "https://horahora-schemas.org/configs/environment-variables.schema.json",
  title: "EnvironmentVariables",
  description: "The environment variables used by Horahora project.",
  type: "object",
  required: ["JWT_KEYPAIR"],
  properties: {
    COMPOSE_PROJECT_NAME: {
      type: "string",
      default: "horahora",
      description:
        "https://docs.docker.com/compose/reference/envvars/#compose_project_name",
    },
    COMPOSE_FILE: {
      type: "string",
      default: "./docker-compose.yml",
      description:
        "https://docs.docker.com/compose/reference/envvars/#compose_file",
    },
    COMPOSE_DOCKER_CLI_BUILD: {
      type: "integer",
      enum: [0, 1],
      default: 1,
      description:
        "https://docs.docker.com/compose/reference/envvars/#compose_docker_cli_build",
    },
    DOCKER_BUILDKIT: {
      type: "integer",
      enum: [0, 1],
      default: 1,
      description: "Enable buildkit for docker.",
    },
    BUCKET_NAME: {
      type: "string",
      default: "otomads",
      description:
        "The storage bucket name for b2 or s3.\nIf you want to use S3, you need to include your aws credentials and config in $HOME/.aws. The config and credentials will be mounted into the relevant services at runtime.\nSee https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html for more information.",
    },
    ORIGIN_FQDN: {
      type: "string",
      default: "http://localhost:9000/$BUCKET_NAME",
      description:
        "For non-local applications, should be the public address of the server (if minio) or storage origin.\nThis will be the public URL of your Backblaze bucket WITH NO TRAILING SLASH. E.g. for me it's: https://f002.backblazeb2.com/file/otomads for backblaze, or https://horahora-dev-otomads.s3-us-west-1.amazonaws.com for s3.",
    },
    STORAGE_BACKEND: {
      type: "string",
      enum: ["b2", "s3"],
      default: "s3",
    },
    STORAGE_API_ID: {
      type: "string",
      default: "minioadmin",
      description:
        "The API ID for your Backblaze account if using backblaze, otherwise blank",
    },
    STORAGE_API_KEY: {
      type: "string",
      default: "minioadmin",
      description: "The API key for your Backblaze account, otherwise blank",
    },
    STORAGE_ENDPOINT: {
      type: "string",
      default: "http://minio:9000",
    },
    VPN_PROVIDER: {
      type: "string",
      default: "mullvad",
    },
    VPN_IP: {
      type: "string",
      default: "10.64.168.1/32",
    },
    WIREGUARD_PRIVATE_KEY: {
      type: "string",
    },
    VPN_CITIES: {
      type: "string",
      default: "Tokyo",
    },
    SOCKS_ADDR: {
      type: "string",
    },
    GOLANG_PROTOBUF_REGISTRATION_CONFLICT: {
      type: "string",
      default: "warn",
    },
    GRPCPort: {
      type: "integer",
      default: 7777,
    },
    JWT_KEYPAIR: {
      type: "string",
    },
  },
} as const;
export const archivalEventSchema = {
  $id: "https://horahora-schemas.org/public-api/archival-event.schema.json",
  title: "ArchivalEvent",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    video_url: {
      type: "string",
    },
    parent_url: {
      type: "string",
    },
    message: {
      type: "string",
    },
    timestamp: {
      type: "string",
    },
  },
} as const;
export const archivalRequestSchema = {
  $id: "https://horahora-schemas.org/public-api/archival-request.schema.json",
  title: "ArchivalRequest",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    UserID: {
      type: "integer",
    },
    Url: {
      type: "string",
    },
    ArchivedVideos: {
      type: "integer",
    },
    CurrentTotalVideos: {
      type: "integer",
    },
    LastSynced: {
      type: "string",
    },
    BackoffFactor: {
      type: "integer",
    },
    DownloadID: {
      type: "integer",
    },
    UndownloadableVideos: {
      type: "integer",
    },
  },
} as const;
export const auditDataSchema = {
  $id: "https://horahora-schemas.org/public-api/audit-data.schema.json",
  title: "AuditData",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    Length: {
      type: "integer",
    },
    Events: {
      type: "array",
      items: {
        $ref: "https://horahora-schemas.org/public-api/audit-event.schema.json",
      },
    },
  },
} as const;
export const auditEventSchema = {
  $id: "https://horahora-schemas.org/public-api/audit-event.schema.json",
  title: "AuditEvent",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    ID: {
      type: "integer",
    },
    UserID: {
      type: "integer",
    },
    Message: {
      type: "string",
    },
    Timestamp: {
      type: "string",
    },
  },
} as const;
export const commentDataSchema = {
  $id: "https://horahora-schemas.org/public-api/comment-data.schema.json",
  title: "CommentData",
  type: "object",
  required: [
    "id",
    "created",
    "content",
    "fullname",
    "profile_picture_url",
    "upvote_count",
    "user_has_upvoted",
  ],
  additionalProperties: false,
  properties: {
    id: {
      type: "integer",
    },
    created: {
      type: "string",
    },
    content: {
      type: "string",
    },
    fullname: {
      type: "string",
    },
    profile_picture_url: {
      type: "string",
    },
    upvote_count: {
      type: "integer",
    },
    user_has_upvoted: {
      type: "boolean",
    },
    parent: {
      type: "integer",
    },
  },
} as const;
export const loggedInUserDataSchema = {
  $id: "https://horahora-schemas.org/public-api/current-user-profile.schema.json",
  title: "LoggedInUserData",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    userID: {
      type: "integer",
    },
    username: {
      type: "string",
    },
    profile_picture_url: {
      type: "string",
    },
    rank: {
      type: "integer",
    },
    banned: {
      type: "boolean",
    },
  },
} as const;
export const paginationDataSchema = {
  $id: "https://horahora-schemas.org/public-api/pagination-data.schema.json",
  title: "PaginationData",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    NumberOfItems: {
      type: "integer",
    },
    CurrentPage: {
      type: "integer",
    },
  },
} as const;
export const profileDataSchema = {
  $id: "https://horahora-schemas.org/public-api/profile-data.schema.json",
  title: "ProfileData",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    PaginationData: {
      $ref: "https://horahora-schemas.org/public-api/pagination-data.schema.json",
    },
    UserID: {
      type: "integer",
    },
    Username: {
      type: "string",
    },
    ProfilePictureURL: {
      type: "string",
    },
    Videos: {
      type: "array",
      items: {
        $ref: "https://horahora-schemas.org/public-api/video.schema.json",
      },
    },
    Banned: {
      type: "boolean",
    },
    L: {
      $ref: "https://horahora-schemas.org/public-api/current-user-profile.schema.json",
    },
  },
} as const;
export const fileUploadSchema = {
  $id: "https://horahora-schemas.org/public-api/upload.schema.json",
  title: "FileUpload",
  description: "Schema for the upload form",
  required: ["file[0]", "file[1]"],
  additionalProperties: false,
  properties: {
    title: {
      type: "string",
    },
    description: {
      type: "string",
    },
    tags: {
      type: "array",
      items: {
        type: "string",
      },
    },
    "file[0]": {
      type: "string",
      description: "Video file.",
    },
    "file[1]": {
      type: "string",
      description: "Thumbnail file.",
    },
  },
} as const;
export const videoDetailSchema = {
  $id: "https://horahora-schemas.org/public-api/video-detail.schema.json",
  title: "VideoDetail",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    Title: {
      type: "string",
    },
    MPDLoc: {
      type: "string",
    },
    Views: {
      type: "integer",
    },
    Rating: {
      type: "number",
    },
    VideoID: {
      type: "integer",
    },
    AuthorID: {
      type: "integer",
    },
    Username: {
      type: "string",
    },
    UserDescription: {
      type: "string",
    },
    VideoDescription: {
      type: "string",
    },
    UserSubscribers: {
      type: "integer",
    },
    ProfilePicture: {
      type: "string",
    },
    UploadDate: {
      type: "string",
      description: "should be a datetime",
    },
    Tags: {
      type: "array",
      items: {
        type: "string",
      },
    },
    RecommendedVideos: {
      type: "array",
      items: {
        $ref: "https://horahora-schemas.org/public-api/video.schema.json",
      },
    },
    L: {
      $ref: "https://horahora-schemas.org/public-api/current-user-profile.schema.json",
    },
  },
} as const;
export const videoInProgressSchema = {
  $id: "https://horahora-schemas.org/public-api/video-in-progress.schema.json",
  title: "VideoInProgress",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    Website: {
      type: "string",
    },
    VideoID: {
      type: "string",
    },
    DlStatus: {
      type: "string",
    },
  },
} as const;
export const videoMetadataSchema = {
  $id: "https://horahora-schemas.org/public-api/video-metadata.schema.json",
  title: "VideoMetadata",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    videoLoc: {
      type: "string",
      description: "The location of the DASH manifest",
    },
    videoTitle: {
      type: "string",
    },
    rating: {
      type: "number",
    },
    authorName: {
      type: "string",
      description: "Do I need this? probably not",
    },
    views: {
      type: "integer",
    },
    videoID: {
      type: "integer",
    },
    uploadDate: {
      type: "string",
    },
    description: {
      type: "string",
    },
    authorID: {
      type: "integer",
    },
    tags: {
      type: "array",
      items: {
        type: "string",
      },
    },
  },
} as const;
export const videoSchema = {
  $id: "https://horahora-schemas.org/public-api/video.schema.json",
  title: "Video",
  type: "object",
  required: [],
  additionalProperties: false,
  properties: {
    Title: {
      type: "string",
    },
    VideoID: {
      type: "integer",
    },
    Views: {
      type: "integer",
    },
    AuthorID: {
      type: "integer",
    },
    AuthorName: {
      type: "string",
    },
    ThumbnailLoc: {
      type: "string",
    },
    Rating: {
      type: "number",
    },
  },
} as const;
export const schemaMap = {
  [environmentVariablesSchema.$id]: environmentVariablesSchema,
  [archivalEventSchema.$id]: archivalEventSchema,
  [archivalRequestSchema.$id]: archivalRequestSchema,
  [auditDataSchema.$id]: auditDataSchema,
  [auditEventSchema.$id]: auditEventSchema,
  [commentDataSchema.$id]: commentDataSchema,
  [loggedInUserDataSchema.$id]: loggedInUserDataSchema,
  [paginationDataSchema.$id]: paginationDataSchema,
  [profileDataSchema.$id]: profileDataSchema,
  [fileUploadSchema.$id]: fileUploadSchema,
  [videoDetailSchema.$id]: videoDetailSchema,
  [videoInProgressSchema.$id]: videoInProgressSchema,
  [videoMetadataSchema.$id]: videoMetadataSchema,
  [videoSchema.$id]: videoSchema,
} as const;
