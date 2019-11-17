package renewable

import (
	"testing"

	"github.com/SladeThe/common-go/renewable/periods"
)

func TestSoftHard_Get(t *testing.T) {
	simpleTestRenewable(t, func(produce ProduceFunc) Renewable {
		softHard := periods.Periods{Default: defaultPeriod, Error: errorPeriod}
		return SoftHardNoCtx(softHard, softHard, produce)
	})
}

func TestSoftHard_Get_Async(t *testing.T) {
	asyncTestRenewable(t, func(produce ProduceFunc) Renewable {
		softHard := periods.Periods{Default: defaultPeriod, Error: errorPeriod}
		return SoftHardNoCtx(softHard, softHard, produce)
	}, asyncTestRenewableOnce)
}

func TestSoftHard_Get_AsyncSoftHard(t *testing.T) {
	asyncTestRenewable(t, func(produce ProduceFunc) Renewable {
		return SoftHardNoCtx(
			periods.Periods{Default: defaultPeriod, Error: errorPeriod},
			periods.Periods{Default: defaultPeriod * 2, Error: errorPeriod * 2},
			produce,
		)
	}, asyncTestRenewableOnce)
}
