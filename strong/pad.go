package strong

import (
	"fmt"
	"strings"
)

func Pad(str string, width int) string {
	if padding := width - len(str); padding > 0 {
		return fmt.Sprintf("%s%s", strings.Repeat(" ", padding), str)
	}
	return str
}
