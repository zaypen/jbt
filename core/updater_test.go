package core

import (
	"github.com/zaypen/gest"
	"testing"
)

func TestCompareVersion(t *testing.T) {
	gest.I(t).Should("1.0 > 0.9").Expect(compareVersion("1.0", "0.9") > 0).ToBe(true)
	gest.I(t).Should("9.1 < 2018.1").Expect(compareVersion("9.1", "2018.1") < 0).ToBe(true)
	gest.I(t).Should("2018.1.4 < 2018.1.5").Expect(compareVersion("2018.1.4", "2018.1.5") < 0).ToBe(true)
	gest.I(t).Should("2018.2 < 2018.2.1").Expect(compareVersion("2018.2", "2018.2.1") < 0).ToBe(true)
}
