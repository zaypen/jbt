package updater

import (
	"testing"
	"github.com/zaypen/jbt/test"
)

func TestCompareVersion(t *testing.T)  {
	test.It(t).Should("1.0 > 0.9").Expect(compareVersion("1.0", "0.9") > 0).ToBe(true)
	test.It(t).Should("0.9.0 == 0.9").Expect(compareVersion("0.9.0", "0.9") == 0).ToBe(true)
	test.It(t).Should("2018.1.4 < 2018.1.5").Expect(compareVersion("2018.1.4", "2018.1.5") < 0).ToBe(true)
}
