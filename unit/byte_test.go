package unit

import (
	"github.com/zaypen/gest"
	"testing"
)

func TestBytes(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		expected string
	}{
		{"Bytes(100)=100B", 100, "100B"},
		{"Bytes(3620)=3.6KB", 3620, "3.5KB"},
		{"Bytes(3651)=3.6KB", 3651, "3.6KB"},
		{"Bytes(5221589)=5.2MB", 5221589, "5.0MB"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gest.I(t).Should(test.name).Expect(Bytes(test.value)).ToBe(test.expected)
		})
	}
}
