package schemamap

import (
	"fmt"
	"horahora/cli/src/lib/codegen"
	"horahora/cli/src/lib/json"
	"horahora/cli/src/lib/strings"
	"path/filepath"
)

const packageName = "schemamap"

var CodeGenerator codegen.ICodegenFunc = codegen.CreateCodegenModule(generateSchemaMap)

func generateSchemaMap() (string, string) {
	genPath := filepath.Join(codegen.CodegenFolderPath, "schema", "map")
	var varDeclarations []string

	schemaCollection := codegen.CollectJSONSchemas()

	for _, jsonContent := range *schemaCollection {
		schema := json.FromJSON[codegen.IJSONSchema](jsonContent)
		schemaTitle := schema["title"].(string)

		varDeclaration := strings.MultilineString(
			fmt.Sprintf("var %v []byte = []byte(`%v`)", schemaTitle, string(jsonContent)),
		)

		varDeclarations = append(varDeclarations, varDeclaration)
	}

	codeContent := strings.MultilineString(
		fmt.Sprintf("package %v", packageName),
		strings.MultilineString(varDeclarations...),
	)

	return codeContent, genPath
}
