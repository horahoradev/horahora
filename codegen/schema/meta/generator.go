package meta

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const schemaFolder string = "schema"
const metaSchemaFileName = "meta.schema.json"
const resultFilename = "export.go"

// Inline the metaschema JSON into a golang module.
//
// TODO: generate interface/struct of said JSON.
func GenerateMetaSchema() {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	metaSchemaPath := filepath.Join(cwd, "..", schemaFolder, metaSchemaFileName)
	metaSchemaJSON, err := os.ReadFile(metaSchemaPath)

	if err != nil {
		panic(err)
	}

	outputPath := filepath.Join(cwd, "schema", "meta", resultFilename)
	outputSlice := []string{
		"package meta",
		fmt.Sprintf("const MetaJSONSchema string = `%v`", string(metaSchemaJSON)),
	}
	outputContent := strings.Join(outputSlice, "\n")

	writeErr := os.WriteFile(outputPath, []byte(outputContent), fs.ModePerm)

	if writeErr != nil {
		panic(err)
	}
}
