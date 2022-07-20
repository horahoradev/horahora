package schema

import (
	schemamap "horahora/cli/src/codegen/schema/map"
)

func GenerateSchemas() {
	schemamap.CodeGenerator()
}
