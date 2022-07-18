/*
  Codegen is only run in development,
  therefore a nodejs server context.
*/

import { runCodegen } from "./lib";

(async () => {
  await runCodegen()
})()
