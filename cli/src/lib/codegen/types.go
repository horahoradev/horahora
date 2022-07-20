package codegen

// Generic representation of parsed JSON schema.
type IJSONSchema = map[string]any

// Key is the schema ID while the value is the schema literal.
type ISchemaCollection = map[string][]byte

// A function which creates a code string
// and returns it along the folder it has to be generated at.
type IGeneratorFunc = func() (string, string)
type ICodegenFunc = func()
