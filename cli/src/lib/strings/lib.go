package strings

import (
	"fmt"
	"strings"
)

// Creates a multiline comment string out of provided string arguments
func CommentMultiline(lines ...string) string {
	outputSlice := []string{"/*"}

	for _, line := range lines {
		outputSlice = append(outputSlice, fmt.Sprintf(" * %v", line))
	}

	outputSlice = append(outputSlice, " */")

	return strings.Join(outputSlice, "\n")
}
