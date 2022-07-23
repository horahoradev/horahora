package codegen

import (
	"horahora/cli/src/codegen/schema"
)

func RunCodegen() {
	schema.GenerateSchemas()
}
