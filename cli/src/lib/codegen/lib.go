package codegen

import (
	"fmt"
	"horahora/cli/src/lib/errors"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	schemaFolder      string = "schema"
	schemaFilenameEnd string = ".schema.json"
)

var (
	metaSchemaFilename = fmt.Sprintf("meta%v", schemaFilenameEnd)
	schemaFolderPath   string
)

func CollectJSONSchemas() {
	filepath.WalkDir(schemaFolderPath, func(path string, enry fs.DirEntry, err error) error {
		return nil
	})
}

func init() {
	cwd, err := os.Getwd()
	errors.CheckError(err)
	schemaFolderPath = filepath.Join(cwd, schemaFolder)
}
