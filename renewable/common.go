package renewable

import (
	"context"
	"time"

	"github.com/SladeThe/common-go/renewable/periods"
)

type result struct {
	val  interface{}
	err  error
	time time.Time
}

func (r result) isValidAt(moment time.Time, periods periods.Periods) bool {
	return !r.time.IsZero() && r.time.Add(periods.Period(r.err)).After(moment)
}

func (r result) isValidNow(periods periods.Periods) bool {
	return r.isValidAt(time.Now(), periods)
}

type producing struct {
	ctx     context.Context
	produce ProduceFunc
}
