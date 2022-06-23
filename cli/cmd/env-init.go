package cmd

import (
	"fmt"
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
	Short: "Initialize environment variables file",
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

	if pathErr != nil {
		panic(pathErr)
	}

	envFilePath := filepath.Join(workingFolder, inputEnvFileName)
	envMap := resolveEnvTemplate(envFilePath)
	envMapErr := modifyEnvMap(envMap)

	if envMapErr != nil {
		panic(envMapErr)
	}

	writeComponseConfig()
}

// Analyzes the input env file and returns it as a map.
func resolveEnvTemplate(filePath string) map[string]string {
	file, fileErr := os.Open(filePath)

	if fileErr != nil {
		panic(fileErr)
	}

	defer file.Close()

	envMap := make(map[string]string)

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

func createJWTKeyPair(isDefault bool) string {
	return defaultKeyPairFileName
}

func resolveComposeConfigTemplate() string {
	return inputComposeFileName
}

func writeComponseConfig() string {
	return outputComposeFileName
}
