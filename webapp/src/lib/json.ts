export function toJSON<InputType = unknown>(inputValue: InputType) {
  return JSON.stringify(inputValue);
}
export function fromJSON<OutputType = unknown>(inputJSON: string) {
  return JSON.parse(inputJSON) as OutputType;
}
