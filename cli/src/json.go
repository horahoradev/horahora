package lib

import "encoding/json"

// Parses json string and returns the value of provided type.
func ParseJSON[outputType any](inputJSON string) outputType {
	jsonContent := []byte(inputJSON)
	var jsonResult outputType
	err := json.Unmarshal(jsonContent, &jsonResult)
	CheckError(err)

	return jsonResult
}

// Turns a provided value into a json string.
func StringifyJSON[inputType any](inputItem inputType) string {
	jsonString, err := json.Marshal(inputItem)
	CheckError(err)

	return string(jsonString)
}

// Turns a provided value into a human-readable json string.
func PrettyJSON[inputType any](inputItem inputType) string {
	prettyJSONString, err := json.MarshalIndent(inputItem, "", "  ")
	CheckError(err)

	return string(prettyJSONString)
}
