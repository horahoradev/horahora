package codegen

import (
	"fmt"
	"horahora/cli/src/lib/errors"
	fslib "horahora/cli/src/lib/fs"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Key is the schema ID while the value is the schema literal
type ISchemaCollection = map[string][]byte

const (
	schemaFolder      string = "schema"
	schemaFilenameEnd string = ".schema.json"
)

var (
	schemaCollection   ISchemaCollection
	metaSchemaFilename = fmt.Sprintf("meta%v", schemaFilenameEnd)
	schemaFolderPath   string
)

func CollectJSONSchemas() *ISchemaCollection {
	if len(schemaCollection) != 0 {
		return &schemaCollection
	}

	filepath.WalkDir(schemaFolderPath, func(path string, entry fs.DirEntry, err error) error {
		if !isSchemaFile(entry) {
			return nil
		}

		fslib.ReadFile(path)

		return nil
	})

	return &schemaCollection
}

func isSchemaFile(entry fs.DirEntry) bool {
	// exclusing metaschema from the collection for now
	return entry.Type().IsRegular() && entry.Name() != metaSchemaFilename && strings.HasSuffix(entry.Name(), schemaFilenameEnd)
}

func init() {
	cwd, err := os.Getwd()
	errors.CheckError(err)
	schemaFolderPath = filepath.Join(cwd, schemaFolder)
}
