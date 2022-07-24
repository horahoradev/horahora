import path from "node:path";
import { cwd } from "node:process";

import {
  collectJSONSchemas,
  type ICodegen,
  type ICodegenFunc,
} from "#server/codegen";
import { decapitalizeString, multilineString } from "#lib/strings";

export {
  archivalEventSchema,
  archivalRequestSchema,
  auditDataSchema,
  auditEventSchema,
  commentDataSchema,
  environmentVariablesSchema,
  fileUploadSchema,
  paginationDataSchema,
  profileDataSchema,
  videoDetailSchema,
  schemaMap,
  videoInProgressSchema,
  videoMetadataSchema,
  videoSchema,
  accountSchema,
  accountInitSchema,
  accountClientSchema,
} from "../000_map/_index";
export type {
  IArchivalEvent,
  IArchivalRequest,
  IAuditData,
  IAuditEvent,
  ICommentData,
  IEnvironmentVariables,
  IFileUpload,
  IPaginationData,
  IProfileData,
  IVideo,
  IVideoDetail,
  IVideoInProgress,
  IVideoMetadata,
  IAccount,
  IAccountClient,
  IAccountInit,
} from "../001_interfaces/_index";
export { createAJV, createValidator } from "#lib/json-schema";
export type { ISchemaMap, IJSONSchema } from "#lib/json-schema";

const schemaFolder = ["..", "schema"];

const generateJSONSchemas: ICodegenFunc = async () => {
  const schemaFolderPath = path.join(cwd(), ...schemaFolder);
  const schemaCollection = await collectJSONSchemas(schemaFolderPath);
  const typeImports: string[] = ["ISchemaMap", "IJSONSchema"];
  const results: string[] = [
    "const ajv = createAJV(schemaMap as unknown as ISchemaMap)",
  ];
  const concreteExports: string[] = [];

  for await (const [schemaID, schema] of Object.entries(schemaCollection)) {
    const funcName = `validate${schema.title}`;
    const interfaceName = `I${schema.title}`;
    const schemaName = `${decapitalizeString(schema.title)}Schema`;
    const validatorFunc = `export const ${funcName} = createValidator<${interfaceName}>(${schemaName} as unknown as IJSONSchema, ajv);`;

    typeImports.push(interfaceName);
    concreteExports.push(funcName);
    results.push(validatorFunc);
  }

  const typeImportsString = `import type {${typeImports.join(
    ", "
  )}} from "./generator"`;

  const codegen: ICodegen = {
    exports: {
      concrete: concreteExports,
    },
    result: multilineString(typeImportsString, ...results),
  };

  return codegen;
};

export default generateJSONSchemas;
