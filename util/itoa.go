package util

import "strconv"

func Atoi(a string) int {
	if i, err := strconv.Atoi(a); err == nil {
		return i
	}
	return 0
}
