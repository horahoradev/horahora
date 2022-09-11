package env

import (
	"fmt"
	"horahora/cli/src/lib/errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"github.com/spf13/cobra"
)

// `env validate`
var envValidate = &cobra.Command{
	Use:   "validate",
	Short: "Validate environment variables",
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting environment variables validation...")
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finished environment variables validation.")
	},
	Run: validateEnvVars,
}

// Compares example and present env.
func validateEnvVars(cmd *cobra.Command, args []string) {
	workingFolder, pathErr := os.Getwd()
	errors.CheckError(pathErr)

	envFilePath := filepath.Join(workingFolder, envFileName)
	envMap, fileErr := godotenv.Read(envFilePath)
	errors.CheckError(fileErr)

	exampleENVFilepath := filepath.Join(workingFolder, configsFolder, envExampleFilename)
	exampleENVMap, fileErr := godotenv.Read(exampleENVFilepath)
	errors.CheckError(fileErr)

	missingKeys := compareEnvMaps(envMap, exampleENVMap)

	if len(missingKeys) > 0 {
		keyList := strings.Join(missingKeys, "\n")
		println("\".env\" misses these keys:", keyList)
	}
}

// Compares the present env map to the example one
// and returns the list of missing keys in the former.
//
// @TODO: compare against schema instead
func compareEnvMaps(envMap, exampleENVMap map[string]string) []string {
	var missingKeys []string

	for key := range exampleENVMap {
		_, ok := envMap[key]

		if !ok {
			missingKeys = append(missingKeys, key)
		}
	}

	return missingKeys
}
