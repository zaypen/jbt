package updater

import (
	"testing"
	"github.com/zaypen/gest"
)

func TestCompareVersion(t *testing.T)  {
	gest.I(t).Should("1.0 > 0.9").Expect(compareVersion("1.0", "0.9") > 0).ToBe(true)
	gest.I(t).Should("0.9.0 == 0.9").Expect(compareVersion("0.9.0", "0.9") == 0).ToBe(true)
	gest.I(t).Should("2018.1.4 < 2018.1.5").Expect(compareVersion("2018.1.4", "2018.1.5") < 0).ToBe(true)
}
