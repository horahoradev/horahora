package cmd

import (
	"fmt"
	"horahora/cli/src/lib/errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

const inputEnvFileName string = "secrets.env.template"
const inputComposeFileName string = "docker-compose.yml.envs"
const outputComposeFileName string = "docker-compose.yml"
const defaultKeyPairFileName string = "default_keypair.pem"

// envInitCMD represents the `env init` command
var envInitCMD = &cobra.Command{
	Use:   "init",
	Short: "Initialize environment variables",
	Long:  `Environment variables initiaizer for Horahora.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing environment variables...")
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finished environment variables initialization.")
	},
	Run: initEnvVars,
}

// Initializes `docker-compose` config off files.
func initEnvVars(cmd *cobra.Command, args []string) {
	workingFolder, pathErr := os.Getwd()
	errors.CheckError(pathErr)

	envFilePath := filepath.Join(workingFolder, inputEnvFileName)
	envMap := resolveEnvTemplate(envFilePath)

	composeContent := resolveComposeConfigTemplate(envMap)
	writeComponseConfig(composeContent)
}

// Analyzes the input env file and returns it as a map.
func resolveEnvTemplate(filePath string) map[string]string {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		panic(fileErr)
	}

	defer file.Close()

	envMap := make(map[string]string)

	envMapErr := modifyEnvMap(envMap)
	errors.CheckError(envMapErr)

	return envMap
}

// Applies conditional changes to the values of some keys.
func modifyEnvMap(envMap map[string]string) error {
	useVPNConfig, err := strconv.ParseBool(envMap["USE_VPN_CONFIG"])

	if !useVPNConfig || err != nil {
		envMap["SOCKS_ADDR"] = ""
	}

	useCustomKeypair, err2 := strconv.ParseBool(envMap["USE_VPN_CONFIG"])

	if !useCustomKeypair || err2 != nil {
		envMap["JWT_KEYPAIR"] = createJWTKeyPair(true)
	} else {
		envMap["JWT_KEYPAIR"] = createJWTKeyPair(false)
	}

	buildLocalImages, err3 := strconv.ParseBool(envMap["BUILD_LOCAL_IMAGES"])

	if !buildLocalImages || err3 != nil {
		envMap["FRONTEND_BUILD_DIRECTIVE"] = "image: ghcr.io/horahoradev/horahora:master_frontend"
		envMap["VIDEOSERVICE_BUILD_DIRECTIVE"] = "image: ghcr.io/horahoradev/horahora:master_videoservice"
		envMap["USERSERVICE_BUILD_DIRECTIVE"] = "image: ghcr.io/horahoradev/horahora:master_userservice"
		envMap["SCHEDULER_BUILD_DIRECTIVE"] = "image: ghcr.io/horahoradev/horahora:master_scheduler"
	}

	return nil
}

// Creates jwt key or returns a default one.
func createJWTKeyPair(isDefault bool) string {
	return defaultKeyPairFileName
}

// Resolves the template of compsoe config with env variables.
func resolveComposeConfigTemplate(envMap map[string]string) string {
	return inputComposeFileName
}

// Saves the final compose config or updates the existing one.
func writeComponseConfig(composeContent string) string {
	return outputComposeFileName
}
