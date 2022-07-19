/*
  This module was created by the codegen, do not edit it manually.
*/
import {
  archivalEventSchema,
  archivalRequestSchema,
  auditDataSchema,
  auditEventSchema,
  commentDataSchema,
  environmentVariablesSchema,
  fileUploadSchema,
  loggedInUserDataSchema,
  paginationDataSchema,
  profileDataSchema,
  videoDetailSchema,
  schemaMap,
  videoInProgressSchema,
  videoMetadataSchema,
  videoSchema,
  createAJV,
  createValidator,
} from "./generator";
import type {
  ISchemaMap,
  IJSONSchema,
  IEnvironmentVariables,
  IArchivalEvent,
  IArchivalRequest,
  IAuditData,
  IAuditEvent,
  ICommentData,
  ILoggedInUserData,
  IPaginationData,
  IProfileData,
  IFileUpload,
  IVideoDetail,
  IVideoInProgress,
  IVideoMetadata,
  IVideo,
} from "./generator";
const ajv = createAJV(schemaMap as unknown as ISchemaMap);
export const validateEnvironmentVariables =
  createValidator<IEnvironmentVariables>(
    environmentVariablesSchema as unknown as IJSONSchema,
    ajv
  );
export const validateArchivalEvent = createValidator<IArchivalEvent>(
  archivalEventSchema as unknown as IJSONSchema,
  ajv
);
export const validateArchivalRequest = createValidator<IArchivalRequest>(
  archivalRequestSchema as unknown as IJSONSchema,
  ajv
);
export const validateAuditData = createValidator<IAuditData>(
  auditDataSchema as unknown as IJSONSchema,
  ajv
);
export const validateAuditEvent = createValidator<IAuditEvent>(
  auditEventSchema as unknown as IJSONSchema,
  ajv
);
export const validateCommentData = createValidator<ICommentData>(
  commentDataSchema as unknown as IJSONSchema,
  ajv
);
export const validateLoggedInUserData = createValidator<ILoggedInUserData>(
  loggedInUserDataSchema as unknown as IJSONSchema,
  ajv
);
export const validatePaginationData = createValidator<IPaginationData>(
  paginationDataSchema as unknown as IJSONSchema,
  ajv
);
export const validateProfileData = createValidator<IProfileData>(
  profileDataSchema as unknown as IJSONSchema,
  ajv
);
export const validateFileUpload = createValidator<IFileUpload>(
  fileUploadSchema as unknown as IJSONSchema,
  ajv
);
export const validateVideoDetail = createValidator<IVideoDetail>(
  videoDetailSchema as unknown as IJSONSchema,
  ajv
);
export const validateVideoInProgress = createValidator<IVideoInProgress>(
  videoInProgressSchema as unknown as IJSONSchema,
  ajv
);
export const validateVideoMetadata = createValidator<IVideoMetadata>(
  videoMetadataSchema as unknown as IJSONSchema,
  ajv
);
export const validateVideo = createValidator<IVideo>(
  videoSchema as unknown as IJSONSchema,
  ajv
);
