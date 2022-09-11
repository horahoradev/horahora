package json

import (
	"encoding/json"
	"horahora/cli/src/lib/errors"
)

// Parses json bytes and returns the value of provided type.
func FromJSON[outputType any](inputJSON []byte) outputType {
	jsonContent := []byte(inputJSON)
	var jsonResult outputType
	err := json.Unmarshal(jsonContent, &jsonResult)
	errors.CheckError(err)

	return jsonResult
}

// Turns the value into a json bytes.
func ToJSON[inputType any](inputItem inputType) []byte {
	jsonString, err := json.Marshal(inputItem)
	errors.CheckError(err)

	return jsonString
}

// Turns the value into a human-readable json string.
func ToJSONPretty[inputType any](inputItem inputType) string {
	prettyJSONString, err := json.MarshalIndent(inputItem, "", "  ")
	errors.CheckError(err)

	return string(prettyJSONString)
}
