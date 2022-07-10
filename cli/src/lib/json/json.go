package json

import (
	"encoding/json"
	"horahora/cli/src/lib/errors"
)

// Parses json string and returns the value of provided type.
func FromJSON[outputType any](inputJSON string) outputType {
	jsonContent := []byte(inputJSON)
	var jsonResult outputType
	err := json.Unmarshal(jsonContent, &jsonResult)
	errors.CheckError(err)

	return jsonResult
}

// Turns the value into a json string.
func ToJSON[inputType any](inputItem inputType) string {
	jsonString, err := json.Marshal(inputItem)
	errors.CheckError(err)

	return string(jsonString)
}

// Turns the value into a human-readable json string.
func ToJSONPretty[inputType any](inputItem inputType) string {
	prettyJSONString, err := json.MarshalIndent(inputItem, "", "  ")
	errors.CheckError(err)

	return string(prettyJSONString)
}
