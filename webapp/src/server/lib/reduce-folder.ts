import { readdir } from "node:fs/promises";
import { type Dirent } from "node:fs";
import path, { type ParsedPath } from "node:path";

export type IWalkFunc<InitType> = (
  accumulator: InitType,
  parsedPath: ParsedPath,
  dirEntry: Dirent
) => Promise<typeof accumulator>;

/**
 * Walk from the root folder and call the callback
 * on each file and folder within it,
 * visiting all nested folders.
 * @TODO less aggressive recursion.
 * @param rootFolder A folder to start walking from.
 */
export async function reduceFolder<InitType = unknown>(
  rootFolder: string,
  initValue: InitType,
  walkFunc: IWalkFunc<InitType>
): Promise<typeof initValue> {
  try {
    const rootEntries = await readdir(rootFolder, { withFileTypes: true });
    let accumulator: InitType = initValue;

    for await (const dirEntry of rootEntries) {
      const entryPath = path.join(rootFolder, dirEntry.name);
      const parsedPath = path.parse(entryPath);
      // run callback on a file
      if (dirEntry.isFile()) {
        accumulator = await walkFunc(accumulator, parsedPath, dirEntry);
        continue;
      }

      // run the function itself on a folder
      if (dirEntry.isDirectory()) {
        accumulator = await reduceFolder(entryPath, accumulator, walkFunc);
      }
    }

    return accumulator;
  } catch (error) {
    const isProperError = error instanceof Error;

    if (!isProperError) {
      throw error;
    }

    throw new Error(`Failed to walk folder ${rootFolder}`, { cause: error });
  }
}
