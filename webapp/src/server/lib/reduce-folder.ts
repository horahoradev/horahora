import { readdir } from "node:fs/promises";
import { Dirent } from "node:fs";
import path from "node:path";

export type IWalkFunc<InitType> = (
  accumulator: InitType,
  dirEntry: Dirent
) => Promise<typeof accumulator>;

/**
 * Walk from the root folder and call the callback
 * on each file and folder within it,
 * visiting all nested folders.
 * @TODO less aggressive recursion.
 */
export async function reduceFolder<InitType = unknown>(
  rootFolder: string,
  initValue: InitType,
  walkFunc: IWalkFunc<InitType>
): Promise<typeof initValue> {
  const rootEntries = await readdir(rootFolder, { withFileTypes: true });
  let accumulator: InitType = initValue;

  for await (const dirEntry of rootEntries) {
    // run callback on a file
    if (dirEntry.isFile()) {
      accumulator = await walkFunc(accumulator, dirEntry);
      continue;
    }

    // run the function itself on a folder
    if (dirEntry.isDirectory()) {
      const newRootFolder = path.join(rootFolder, dirEntry.name);
      accumulator = await reduceFolder(newRootFolder, accumulator, walkFunc);
    }
  }

  return accumulator;
}
