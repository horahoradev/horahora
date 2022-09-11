package env

import (
	"github.com/spf13/cobra"
)

const (
	envFileName        string = ".env"
	configsFolder      string = "configs"
	envExampleFilename string = ".env.example"
	schemaFolder       string = "schema"
)

// `env`
var ENVCommand = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables",
	Long:  `Environment variables manager for Horahora.`,
}

func init() {
	ENVCommand.AddCommand(envInit)
	ENVCommand.AddCommand(envValidate)
}
