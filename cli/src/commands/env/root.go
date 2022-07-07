package env

import (
	"github.com/spf13/cobra"
)

// ENVCommand represents the `env` command
var ENVCommand = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables",
	Long:  `Environment variables manager for Horahora.`,
}

func init() {
	ENVCommand.AddCommand(ENVValidate)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
