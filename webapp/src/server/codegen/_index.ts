/*
  Codegen is only run in development,
  therefore a nodejs server context.
*/
export { runCodegen } from "./lib";
export { collectJSONSchemas, createInterfaceFromSchema } from "./json-schema";
export type { IJSONSchemaCollection } from "./json-schema";
export type {
  ICodegen,
  ICodegenExport,
  ICodegenModule,
  ICodegenFunc,
} from "./types";
