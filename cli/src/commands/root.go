package commands

import (
	"fmt"
	codegenCommand "horahora/cli/src/commands/codegen"
	envCommand "horahora/cli/src/commands/env"
	"os"

	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "horahora",
		Short: "A CLI tool to manage horahora.",
		Long:  `Horahora management CLI tool.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(envCommand.ENVCommand)
	rootCmd.AddCommand(codegenCommand.CodegenCommand)
}
