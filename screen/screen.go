package screen

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Size() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	if out, err := cmd.Output(); err == nil {
		trimmed := strings.Replace(string(out), "\n", "", -1)
		if outs := strings.Split(trimmed, " "); len(outs) == 2 {
			if width, err := strconv.Atoi(outs[0]); err == nil && width > 0 {
				if height, err := strconv.Atoi(outs[1]); err == nil && height > 0 {
					return width, height
				}
			}
		}
	}
	return 80, 24
}
