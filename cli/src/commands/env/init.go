package env

import (
	"fmt"
	"horahora/cli/src/lib/errors"
	"horahora/cli/src/lib/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// `env init`
var envInit = &cobra.Command{
	Use:   "init",
	Short: "Initialize environment variables",
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initiallizing environment variables...")
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finished environment variables initialization.")
	},
	Run: initEnvVars,
}

// Copies the `configs/.env.example` file into `.env`
//
// @TODO: create `.env` from schema
func initEnvVars(cmd *cobra.Command, args []string) {
	workingFolder, pathErr := os.Getwd()
	errors.CheckError(pathErr)

	envFilePath := filepath.Join(workingFolder, envFileName)

	if fs.IsExist(envFilePath) {
		fmt.Println("\".env\" file already exists, skipping initialization.")
		return
	}

	exampleENVFilepath := filepath.Join(workingFolder, configsFolder, envExampleFilename)

	fs.CopyFile(exampleENVFilepath, envFilePath)
}
