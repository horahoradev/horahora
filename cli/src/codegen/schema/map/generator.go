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
	var imports string = "import ( \"horahora/cli/src/lib/codegen\" )"
	var varDeclarations []string
	var varMap = map[string]string{}

	schemaCollection := codegen.CollectJSONSchemas()

	for schemaID, jsonContent := range *schemaCollection {
		schema := json.FromJSON[codegen.IJSONSchema](jsonContent)
		schemaTitle := schema["title"].(string)

		varDeclaration := strings.MultilineString(
			fmt.Sprintf("var %v []byte = []byte(`%v`)", schemaTitle, string(jsonContent)),
		)

		varDeclarations = append(varDeclarations, varDeclaration)
		varMap[schemaID] = schemaTitle
	}

	var schemaMapLines = []string{
		"var JSONChemaMap = codegen.ISchemaCollection{",
	}

	for schemaID, title := range varMap {
		schemaMapLines = append(schemaMapLines, fmt.Sprintf("\"%v\" : %v,", schemaID, title))
	}

	schemaMapLines = append(schemaMapLines, "}")

	codeContent := strings.MultilineString(
		fmt.Sprintf("package %v", packageName),
		imports,
		strings.MultilineString(varDeclarations...),
		strings.MultilineString(schemaMapLines...),
	)

	return codeContent, genPath
}
