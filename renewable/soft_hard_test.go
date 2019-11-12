package renewable

import (
	"testing"

	"github.com/SladeThe/common-go/renewable/periods"
)

func TestSoftHard_Get(t *testing.T) {
	simpleTestRenewable(t, func(produce ProduceFunc) Renewable {
		p := periods.Periods{Default: defaultPeriod, Error: errorPeriod}
		return SoftHardNoCtx(p, p, produce)
	})
}

func TestSoftHard_Get_Async(t *testing.T) {
	for i := 0; i < 10; i++ {
		asyncTestRenewable(t, func(produce ProduceFunc) Renewable {
			p := periods.Periods{Default: defaultPeriod, Error: errorPeriod}
			return SoftHardNoCtx(p, p, produce)
		})

		if t.Failed() {
			t.FailNow()
		}
	}
}
