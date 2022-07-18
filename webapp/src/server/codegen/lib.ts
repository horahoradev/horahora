import path from "node:path";
import fs from "node:fs/promises";
import { cwd } from "node:process";

import prettier from "prettier";

import { type ICodegen, type ICodegenModule } from "./types";

import { reduceFolder } from "#server/lib";
import { multilineString } from "#lib/strings";
import { isError } from "#lib/errors";

const codegenFolder = ["src", "codegen"];
const generatorFilename = "generator";
const generatorModule = `${generatorFilename}.ts`;
const resultFilename = "result";
const resultModule = `${resultFilename}.ts`;
const indexFilename = "_index";
const indexModule = `${indexFilename}.ts`;

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
        const { default: codegen, ...libExports }: ICodegenModule =
          await import(modulePath);
        const { exports, result } = await codegen();
        const importKeys = Object.keys(libExports);
        const imports = importKeys.length
          ? multilineString(
              "import {",
              Object.keys(libExports).join(", "),
              `} from "./${generatorFilename}"`
            )
          : undefined;
        const resultPath = path.join(entryPath.dir, resultModule);
        const finalResult = multilineString(
          "/*",
          "  This module was created by the codegen, do not edit it manually.",
          "*/",
          imports,
          result,
          "\n"
        );
        const indexPath = path.join(entryPath.dir, indexModule);
        const indexContent = createIndex(exports);
        let formattedResult;
        let formattedIndex;

        try {
          formattedResult = prettier.format(finalResult);
          formattedIndex = prettier.format(indexContent);
        } catch (error) {
          console.warn(
            `Failed to format the result at "${resultPath}"`,
            "Continue unformatted"
          );
        }

        await fs.writeFile(
          resultPath,
          formattedResult ? formattedResult : finalResult
        );
        await fs.writeFile(
          indexPath,
          formattedIndex ? formattedIndex : indexContent
        );

        codegenDirs.push(path.relative(codegenPath, entryPath.dir));

        return codegenDirs;
      } catch (error) {
        if (!isError(error)) {
          throw error;
        }

        const message = multilineString(
          "Failed to run codegen",
          `Reason: ${error.message}`
        );
        throw new Error(message, { cause: error });
      }
    }
  );

  return codegenDirs;
}

function createIndex(exports: ICodegen["exports"]): string {
  const concreteExports = [
    "export {",
    (exports.concrete || []).join(", "),
    `} from "./${resultFilename}"`,
  ].join(" ");
  const typeExports = [
    "export type {",
    (exports.types || []).join(", "),
    `} from "./${resultFilename}"`,
  ].join(" ");
  const indexContent = multilineString(concreteExports, typeExports);
  return indexContent;
}
