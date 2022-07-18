import path from "node:path";
import fs from "node:fs/promises";
import { cwd } from "node:process";

import prettier from "prettier";
import { message } from "antd";

import { reduceFolder } from "#server/lib";
import { multilineString } from "#lib/strings";
import { isError } from "#lib/errors";

export interface ICodegenModule extends Record<string, unknown> {
  default: () => Promise<string>;
}

const codegenFolder = ["src", "codegen"];
const generatorFilename = "generator";
const generatorModule = `${generatorFilename}.ts`;
const resultFilename = "result.ts";

export async function runCodegen() {
  const codegenPath = path.join(cwd(), ...codegenFolder);

  const codegenDirs = await reduceFolder<string[]>(
    codegenPath,
    [],
    async (codegenDirs, entryPath, dirEntry) => {
      try {
        if (!dirEntry.isFile() || dirEntry.name !== generatorModule) {
          return codegenDirs;
        }

        const modulePath = path.format(entryPath);
        const { default: codegenFunction, ...libExports }: ICodegenModule =
          await import(modulePath);
        const importKeys = Object.keys(libExports)
        const imports = importKeys.length ? multilineString(
          "import {",
          Object.keys(libExports).join(", "),
          `} from "./${generatorFilename}"`
        ) : undefined;
        const resultString = await codegenFunction();
        const resultPath = path.join(entryPath.dir, resultFilename);
        const finalResult = multilineString(
          "/*",
          "  This module was created by the codegen, do not edit it manually.",
          "*/",
          imports,
          resultString,
          "\n"
        );
        let formattedResult;

        try {
          formattedResult = prettier.format(finalResult);
        } catch (error) {
          console.warn(`Failed to format the result at "${resultPath}"`, "Continue unformatted");
        }


        await fs.writeFile(resultPath, finalResult);

        codegenDirs.push(path.relative(codegenPath, entryPath.dir));

        return codegenDirs;
      } catch (error) {
        if (!isError(error)) {
          throw error
        }

        const message = multilineString("Failed to run codegen",
        `Reason: ${error.message}`)
        throw new Error(message, { cause: error })
      }

    }
  );

  return codegenDirs;
}
