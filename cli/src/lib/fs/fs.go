package fs

import (
	errorsLib "errors"
	"fmt"
	"horahora/cli/src/lib/errors"
	"io/fs"
	"os"
)

// Reads the folder at provided location.
func ReadFolder(dirPath string) []fs.DirEntry {
	folderContent, err := os.ReadDir(dirPath)
	errors.CheckError(err)

	return folderContent
}

// Reads the entire file at provided location
// and returns its content as a string.
func ReadFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	errors.CheckError(err)

	return string(data)
}

// Saves the content of the string at provided file location.
func WriteFile(filePath string, content string) {
	fileContent := []byte(content)
	err := os.WriteFile(filePath, fileContent, 0644)
	errors.CheckError(err)
}

// Overwrite the file at destination path
// with the content of the file at sourcePath.
func OverwriteFile(sourcePath, destinationPath string) {
	content := ReadFile(sourcePath)
	WriteFile(destinationPath, content)
}

// Check for existence of the file/folder at path.
func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)

	return !errorsLib.Is(err, fs.ErrExist)
}

// Copies the file from source to destination.
func CopyFile(sourcePath, destinationPath string) {
	if !IsExist(sourcePath) {
		message := fmt.Sprintf("Source path \"%v\" doesn't exist.", sourcePath)
		panic(message)
	}

	if IsExist(destinationPath) {
		message := fmt.Sprintf("Destination path \"%v\"already exists.", destinationPath)
		panic(message)
	}

	content := ReadFile(sourcePath)
	WriteFile(destinationPath, content)
}
