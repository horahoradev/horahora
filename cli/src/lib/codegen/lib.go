package codegen

import (
	"fmt"
	fslib "horahora/cli/src/lib/fs"
	json "horahora/cli/src/lib/json"
	stringslib "horahora/cli/src/lib/strings"
	"io/fs"
	"path/filepath"
	"strings"
)

const (
	schemaFolder      string = "schema"
	schemaFilenameEnd string = ".schema.json"
	resultFilename    string = "export.go"
)

var (
	// The root folder of generated code.
	CodegenFolderPath  string
	schemaCollection   = ISchemaCollection{}
	metaSchemaFilename = fmt.Sprintf("meta%v", schemaFilenameEnd)
	schemaFolderPath   string
	codegenNotice      string = CommentMultiline("This file was created by the codegen, do not edit it manualy.")
)

func CreateCodegenModule(generatorFunc IGeneratorFunc) ICodegenFunc {
	return func() {
		codeContent, genPath := generatorFunc()
		finalContent := stringslib.MultilineString(
			codegenNotice,
			codeContent,
		)
		resultFilePath := filepath.Join(genPath, resultFilename)

		fslib.WriteFile(resultFilePath, []byte(finalContent))
	}
}

// Creates a multiline comment string out of provided string arguments
func CommentMultiline(lines ...string) string {
	var outputSlice []string
	var commentString string

	for _, line := range lines {
		outputSlice = append(outputSlice, fmt.Sprintf(" * %v", line))
	}

	commentString = stringslib.MultilineString(outputSlice...)

	return stringslib.MultilineString(
		"/*",
		commentString,
		" */",
	)
}

// Collect all json schema files into a map.
func CollectJSONSchemas() *ISchemaCollection {
	if len(schemaCollection) != 0 {
		return &schemaCollection
	}

	filepath.WalkDir(schemaFolderPath, func(path string, entry fs.DirEntry, err error) error {
		if !isSchemaFile(entry) {
			return nil
		}

		content := fslib.ReadFile(path)
		schema := json.FromJSON[map[string]any](content)
		schemaID := schema["$id"].(string)

		schemaCollection[schemaID] = content

		return nil
	})

	return &schemaCollection
}

func isSchemaFile(entry fs.DirEntry) bool {
	// excluding metaschema from the collection for now
	return entry.Type().IsRegular() && entry.Name() != metaSchemaFilename && strings.HasSuffix(entry.Name(), schemaFilenameEnd)
}

func init() {
	schemaFolderPath = filepath.Join(fslib.Cwd, schemaFolder)
	CodegenFolderPath = filepath.Join(fslib.Cwd, "cli", "src", "codegen")
}
