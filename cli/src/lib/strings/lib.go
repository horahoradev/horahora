package strings

import (
	"strings"
)

func MultilineString(lines ...string) string {
	return strings.Join(lines, "\n")
}
