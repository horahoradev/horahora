package codegen

import (
	codegenlib "horahora/cli/src/codegen"

	"github.com/spf13/cobra"
)

// `codegen`
var CodegenCommand = &cobra.Command{
	Use:   "codegen",
	Short: "Manage code generation",
	Long:  "Code generation manager for Horahora.",
	Run:   runCodegen,
}

func runCodegen(cmd *cobra.Command, args []string) {
	codegenlib.RunCodegen()
}
