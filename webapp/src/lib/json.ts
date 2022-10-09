export function toJSON<InputType = unknown>(inputValue: InputType): string {
  return JSON.stringify(inputValue);
}

export function fromJSON<OutputType = unknown>(inputJSON: string): OutputType {
  return JSON.parse(inputJSON) as OutputType;
}
