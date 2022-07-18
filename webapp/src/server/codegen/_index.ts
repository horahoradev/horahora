/*
  Codegen is only run in development,
  therefore a nodejs server context.
*/
export type {
  ICodegen,
  ICodegenExport,
  ICodegenModule,
  ICodegenFunc,
} from "./types";
import { runCodegen } from "./lib";

(async () => {
  await runCodegen();
})();
