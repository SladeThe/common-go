package renewable

import (
	"testing"
)

func TestOnDemand_Get(t *testing.T) {
	simpleTestRenewable(t, func(produce ProduceFunc) Renewable {
		return OnDemandErr(period, errPeriod, produce)
	})
}

func TestOnDemand_Get_Async(t *testing.T) {
	for i := 0; i < 10; i++ {
		asyncTestRenewable(t, func(produce ProduceFunc) Renewable {
			return OnDemandErr(period, errPeriod, produce)
		})

		if t.Failed() {
			t.FailNow()
		}
	}
}
