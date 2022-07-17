package codegen

import (
	"horahora/cli/src/codegen/schema/meta"
)

func RunCodegen() {
	meta.GenerateMetaSchema()
}
