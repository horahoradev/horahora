export function multilineString(...lines: Array<string | undefined>) {
  return lines.filter((line) => line).join("\n")
}
