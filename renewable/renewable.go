package renewable

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type ProduceFunc func() (value interface{}, err error)

type Renewable interface {
	// Produce new or get cached value according to the Renewable strategy.
	// It is always thread safe to call this method.
	Get() (interface{}, error)
}

// Returns a new instance of Renewable that will call produce only on demand.
// First time Renewable.Get is called, produce is called and the result is cached.
// Every next time Get is called, it checks the current time against either period or errPeriod,
// depending on the result of the previous produce call.
// If time is expired, produce is called again, otherwise the cached value is returned.
// Neither period nor errPeriod can be negative.
// A zero period is a corner case and means, that produce will be called every time Get is called.
func OnDemandErr(period time.Duration, errPeriod time.Duration, produce ProduceFunc) Renewable {
	if period < 0 {
		panic(fmt.Errorf("period must be zero or positive: %v", period))
	}

	if errPeriod < 0 {
		panic(fmt.Errorf("error period must be zero or positive: %v", errPeriod))
	}

	if produce == nil {
		panic(errors.New("produce must be not nil"))
	}

	return &onDemand{period: period, errPeriod: errPeriod, produce: produce, lock: &sync.RWMutex{}}
}

func OnDemand(period time.Duration, produce ProduceFunc) Renewable {
	return OnDemandErr(period, period, produce)
}

func GetOrPanic(r Renewable) interface{} {
	if value, err := r.Get(); err != nil {
		panic(err)
	} else {
		return value
	}
}
