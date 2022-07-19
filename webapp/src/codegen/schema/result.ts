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

export const environmentVariablesSchema = {
	$id: 'https://horahora-schemas.org/configs/environment-variables.schema.json',
	title: 'EnvironmentVariables',
	description: 'The environment variables used by Horahora project.',
	type: 'object',
	required: [
		'JWT_KEYPAIR'
	],
	properties: {
		COMPOSE_PROJECT_NAME: {
			type: 'string',
			default: 'horahora',
			description: 'https://docs.docker.com/compose/reference/envvars/#compose_project_name'
		},
		COMPOSE_FILE: {
			type: 'string',
			default: './docker-compose.yml',
			description: 'https://docs.docker.com/compose/reference/envvars/#compose_file'
		},
		COMPOSE_DOCKER_CLI_BUILD: {
			type: 'integer',
			enum: [
				0,
				1
			],
			default: 1,
			description: 'https://docs.docker.com/compose/reference/envvars/#compose_docker_cli_build'
		},
		DOCKER_BUILDKIT: {
			type: 'integer',
			enum: [
				0,
				1
			],
			default: 1,
			description: 'Enable buildkit for docker.'
		},
		BUCKET_NAME: {
			type: 'string',
			default: 'otomads',
			description: 'The storage bucket name for b2 or s3.\nIf you want to use S3, you need to include your aws credentials and config in $HOME/.aws. The config and credentials will be mounted into the relevant services at runtime.\nSee https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html for more information.'
		},
		ORIGIN_FQDN: {
			type: 'string',
			default: 'http://localhost:9000/$BUCKET_NAME',
			description: 'For non-local applications, should be the public address of the server (if minio) or storage origin.\nThis will be the public URL of your Backblaze bucket WITH NO TRAILING SLASH. E.g. for me it\'s: https://f002.backblazeb2.com/file/otomads for backblaze, or https://horahora-dev-otomads.s3-us-west-1.amazonaws.com for s3.'
		},
		STORAGE_BACKEND: {
			type: 'string',
			enum: [
				'b2',
				's3'
			],
			default: 's3'
		},
		STORAGE_API_ID: {
			type: 'string',
			default: 'minioadmin',
			description: 'The API ID for your Backblaze account if using backblaze, otherwise blank'
		},
		STORAGE_API_KEY: {
			type: 'string',
			default: 'minioadmin',
			description: 'The API key for your Backblaze account, otherwise blank'
		},
		STORAGE_ENDPOINT: {
			type: 'string',
			default: 'http://minio:9000'
		},
		VPN_PROVIDER: {
			type: 'string',
			default: 'mullvad'
		},
		VPN_IP: {
			type: 'string',
			default: '10.64.168.1/32'
		},
		WIREGUARD_PRIVATE_KEY: {
			type: 'string'
		},
		VPN_CITIES: {
			type: 'string',
			default: 'Tokyo'
		},
		SOCKS_ADDR: {
			type: 'string'
		},
		GOLANG_PROTOBUF_REGISTRATION_CONFLICT: {
			type: 'string',
			default: 'warn'
		},
		GRPCPort: {
			type: 'integer',
			default: 7777
		},
		JWT_KEYPAIR: {
			type: 'string'
		}
	}
} as const
export interface IArchivalEvent {
  video_url?: string;
  parent_url?: string;
  message?: string;
  timestamp?: string;
}

export const archivalEventSchema = {
	$id: 'https://horahora-schemas.org/public-api/archival-event.schema.json',
	title: 'ArchivalEvent',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		video_url: {
			type: 'string'
		},
		parent_url: {
			type: 'string'
		},
		message: {
			type: 'string'
		},
		timestamp: {
			type: 'string'
		}
	}
} as const
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

export const archivalRequestSchema = {
	$id: 'https://horahora-schemas.org/public-api/archival-request.schema.json',
	title: 'ArchivalRequest',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		UserID: {
			type: 'integer'
		},
		Url: {
			type: 'string'
		},
		ArchivedVideos: {
			type: 'integer'
		},
		CurrentTotalVideos: {
			type: 'integer'
		},
		LastSynced: {
			type: 'string'
		},
		BackoffFactor: {
			type: 'integer'
		},
		DownloadID: {
			type: 'integer'
		},
		UndownloadableVideos: {
			type: 'integer'
		}
	}
} as const
export interface IAuditData {
  Length?: number;
  Events?: IAuditEvent[];
}

export const auditDataSchema = {
	$id: 'https://horahora-schemas.org/public-api/audit-data.schema.json',
	title: 'AuditData',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		Length: {
			type: 'integer'
		},
		Events: {
			type: 'array',
			items: {
				$ref: 'https://horahora-schemas.org/public-api/audit-event.schema.json'
			}
		}
	}
} as const
export interface IAuditEvent {
  ID?: number;
  UserID?: number;
  Message?: string;
  Timestamp?: string;
}

export const auditEventSchema = {
	$id: 'https://horahora-schemas.org/public-api/audit-event.schema.json',
	title: 'AuditEvent',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		ID: {
			type: 'integer'
		},
		UserID: {
			type: 'integer'
		},
		Message: {
			type: 'string'
		},
		Timestamp: {
			type: 'string'
		}
	}
} as const
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

export const commentDataSchema = {
	$id: 'https://horahora-schemas.org/public-api/comment-data.schema.json',
	title: 'CommentData',
	type: 'object',
	required: [
		'id',
		'created',
		'content',
		'fullname',
		'profile_picture_url',
		'upvote_count',
		'user_has_upvoted'
	],
	additionalProperties: false,
	properties: {
		id: {
			type: 'integer'
		},
		created: {
			type: 'string'
		},
		content: {
			type: 'string'
		},
		fullname: {
			type: 'string'
		},
		profile_picture_url: {
			type: 'string'
		},
		upvote_count: {
			type: 'integer'
		},
		user_has_upvoted: {
			type: 'boolean'
		},
		parent: {
			type: 'integer'
		}
	}
} as const
export interface ILoggedInUserData {
  userID?: number;
  username?: string;
  profile_picture_url?: string;
  rank?: number;
  banned?: boolean;
}

export const loggedInUserDataSchema = {
	$id: 'https://horahora-schemas.org/public-api/current-user-profile.schema.json',
	title: 'LoggedInUserData',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		userID: {
			type: 'integer'
		},
		username: {
			type: 'string'
		},
		profile_picture_url: {
			type: 'string'
		},
		rank: {
			type: 'integer'
		},
		banned: {
			type: 'boolean'
		}
	}
} as const
export interface IPaginationData {
  NumberOfItems?: number;
  CurrentPage?: number;
}

export const paginationDataSchema = {
	$id: 'https://horahora-schemas.org/public-api/pagination-data.schema.json',
	title: 'PaginationData',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		NumberOfItems: {
			type: 'integer'
		},
		CurrentPage: {
			type: 'integer'
		}
	}
} as const
export interface IProfileData {
  PaginationData?: IPaginationData;
  UserID?: number;
  Username?: string;
  ProfilePictureURL?: string;
  Videos?: IVideo[];
  Banned?: boolean;
  L?: ILoggedInUserData;
}

export const profileDataSchema = {
	$id: 'https://horahora-schemas.org/public-api/profile-data.schema.json',
	title: 'ProfileData',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		PaginationData: {
			$ref: 'https://horahora-schemas.org/public-api/pagination-data.schema.json'
		},
		UserID: {
			type: 'integer'
		},
		Username: {
			type: 'string'
		},
		ProfilePictureURL: {
			type: 'string'
		},
		Videos: {
			type: 'array',
			items: {
				$ref: 'https://horahora-schemas.org/public-api/video.schema.json'
			}
		},
		Banned: {
			type: 'boolean'
		},
		L: {
			$ref: 'https://horahora-schemas.org/public-api/current-user-profile.schema.json'
		}
	}
} as const
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

export const fileUploadSchema = {
	$id: 'https://horahora-schemas.org/public-api/upload.schema.json',
	title: 'FileUpload',
	description: 'Schema for the upload form',
	required: [
		'file[0]',
		'file[1]'
	],
	additionalProperties: false,
	properties: {
		title: {
			type: 'string'
		},
		description: {
			type: 'string'
		},
		tags: {
			type: 'array',
			items: {
				type: 'string'
			}
		},
		'file[0]': {
			type: 'string',
			description: 'Video file.'
		},
		'file[1]': {
			type: 'string',
			description: 'Thumbnail file.'
		}
	}
} as const
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
  L?: ILoggedInUserData;
}

export const videoDetailSchema = {
	$id: 'https://horahora-schemas.org/public-api/video-detail.schema.json',
	title: 'VideoDetail',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		Title: {
			type: 'string'
		},
		MPDLoc: {
			type: 'string'
		},
		Views: {
			type: 'integer'
		},
		Rating: {
			type: 'number'
		},
		VideoID: {
			type: 'integer'
		},
		AuthorID: {
			type: 'integer'
		},
		Username: {
			type: 'string'
		},
		UserDescription: {
			type: 'string'
		},
		VideoDescription: {
			type: 'string'
		},
		UserSubscribers: {
			type: 'integer'
		},
		ProfilePicture: {
			type: 'string'
		},
		UploadDate: {
			type: 'string',
			description: 'should be a datetime'
		},
		Tags: {
			type: 'array',
			items: {
				type: 'string'
			}
		},
		RecommendedVideos: {
			type: 'array',
			items: {
				$ref: 'https://horahora-schemas.org/public-api/video.schema.json'
			}
		},
		L: {
			$ref: 'https://horahora-schemas.org/public-api/current-user-profile.schema.json'
		}
	}
} as const
export interface IVideoInProgress {
  Website?: string;
  VideoID?: string;
  DlStatus?: string;
}

export const videoInProgressSchema = {
	$id: 'https://horahora-schemas.org/public-api/video-in-progress.schema.json',
	title: 'VideoInProgress',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		Website: {
			type: 'string'
		},
		VideoID: {
			type: 'string'
		},
		DlStatus: {
			type: 'string'
		}
	}
} as const
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

export const videoMetadataSchema = {
	$id: 'https://horahora-schemas.org/public-api/video-metadata.schema.json',
	title: 'VideoMetadata',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		videoLoc: {
			type: 'string',
			description: 'The location of the DASH manifest'
		},
		videoTitle: {
			type: 'string'
		},
		rating: {
			type: 'number'
		},
		authorName: {
			type: 'string',
			description: 'Do I need this? probably not'
		},
		views: {
			type: 'integer'
		},
		videoID: {
			type: 'integer'
		},
		uploadDate: {
			type: 'string'
		},
		description: {
			type: 'string'
		},
		authorID: {
			type: 'integer'
		},
		tags: {
			type: 'array',
			items: {
				type: 'string'
			}
		}
	}
} as const
export interface IVideo {
  Title?: string;
  VideoID?: number;
  Views?: number;
  AuthorID?: number;
  AuthorName?: string;
  ThumbnailLoc?: string;
  Rating?: number;
}

export const videoSchema = {
	$id: 'https://horahora-schemas.org/public-api/video.schema.json',
	title: 'Video',
	type: 'object',
	required: [],
	additionalProperties: false,
	properties: {
		Title: {
			type: 'string'
		},
		VideoID: {
			type: 'integer'
		},
		Views: {
			type: 'integer'
		},
		AuthorID: {
			type: 'integer'
		},
		AuthorName: {
			type: 'string'
		},
		ThumbnailLoc: {
			type: 'string'
		},
		Rating: {
			type: 'number'
		}
	}
} as const

