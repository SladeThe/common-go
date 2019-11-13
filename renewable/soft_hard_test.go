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
	asyncTestRenewable(t, func(produce ProduceFunc) Renewable {
		p := periods.Periods{Default: defaultPeriod, Error: errorPeriod}
		return SoftHardNoCtx(p, p, produce)
	})
}

// TODO add more tests
