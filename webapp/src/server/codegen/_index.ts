/*
  Codegen is only run in development,
  therefore a nodejs server context.
*/

import path from "node:path";
import fs from "node:fs/promises";
import { cwd } from "node:process";

import prettier from "prettier";

import { reduceFolder } from "#server/lib";
import { multilineString } from "#lib/strings";

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
      if (!dirEntry.isFile() || dirEntry.name !== generatorModule) {
        return codegenDirs;
      }

      const modulePath = path.format(entryPath);
      const { default: codegenFunction, ...libExports }: ICodegenModule =
        await import(modulePath);
      const imports = multilineString(
        "import {",
        Object.keys(libExports).join(", "),
        `} from "./${generatorFilename}"`
      );
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
      const formattedResult = prettier.format(finalResult);

      await fs.writeFile(resultPath, formattedResult);

      codegenDirs.push(path.relative(codegenPath, entryPath.dir));

      return codegenDirs;
    }
  );

  return codegenDirs;
}
