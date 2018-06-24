package unit

import "fmt"

type unit struct {
	display string
	value   uint64
}

var units = []unit{
	{"TB", 1024 * 1024 * 1024 * 1024},
	{"GB", 1024 * 1024 * 1024},
	{"MB", 1024 * 1024},
	{"KB", 1024},
}

func bytes(n uint64, precision int) string {
	for _, u := range units {
		if n >= u.value {
			format := fmt.Sprintf("%%.%df%%s", precision)
			return fmt.Sprintf(format, float64(n)/float64(u.value), u.display)
		}
	}
	return fmt.Sprintf("%dB", n)
}

func Bytes(n uint64) string {
	return bytes(n, 1)
}
