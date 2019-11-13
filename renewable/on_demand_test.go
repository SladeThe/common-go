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
	asyncTestRenewable(t, func(produce ProduceFunc) Renewable {
		return OnDemandNoCtx(periods.Periods{Default: defaultPeriod, Error: errorPeriod}, produce)
	})
}
