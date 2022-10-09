export function multilineString(...lines: Array<string | undefined>) {
  return lines.filter((line) => line).join("\n");
}

/**
 * Changes the first letter of a string to lower case.
 * @returns Lower-cased string.
 */
export function decapitalizeString(inputString: string) {
  const firstLetter = inputString.charAt(0).toLowerCase();
  const rest = inputString.slice(1);
  const result = [firstLetter, rest].join("");
  return result;
}
