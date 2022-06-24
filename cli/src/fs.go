package lib

import (
	"io/fs"
	"os"
)

// Reads the folder at provided location.
func ReadFolder(dirPath string) []fs.DirEntry {
	folderContent, err := os.ReadDir(dirPath)
	CheckError(err)

	return folderContent
}

// Reads the entire file at provided location
// and returns its content as a string.
func ReadFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	CheckError(err)

	return string(data)
}

// Saves the content of the string at provided file location.
func WriteFile(filePath string, content string) {
	fileContent := []byte(content)
	err := os.WriteFile(filePath, fileContent, 0644)
	CheckError(err)
}
