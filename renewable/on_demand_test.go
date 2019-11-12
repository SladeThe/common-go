package renewable

import (
	"testing"

	"github.com/SladeThe/common-go/renewable/periods"
)

func TestOnDemand_Get(t *testing.T) {
	simpleTestRenewable(t, func(produce ProduceFunc) Renewable {
		return OnDemandNoCtx(periods.Periods{Default: defaultPeriod, Error: errorPeriod}, produce)
	})
}

func TestOnDemand_Get_Async(t *testing.T) {
	for i := 0; i < 10; i++ {
		asyncTestRenewable(t, func(produce ProduceFunc) Renewable {
			return OnDemandNoCtx(periods.Periods{Default: defaultPeriod, Error: errorPeriod}, produce)
		})

		if t.Failed() {
			t.FailNow()
		}
	}
}
