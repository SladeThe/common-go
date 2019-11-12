package renewable

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/SladeThe/common-go/renewable/periods"
)

type ProduceFunc func(ctx context.Context) (value interface{}, err error)

type Renewable interface {
	// Produce new or get cached value according to the Renewable strategy.
	// It is always thread safe to call this method.
	Get() (interface{}, error)
}

// Returns a new instance of Renewable that will call produce only on demand.
// First time Renewable.Get is called, produce is called and the result is cached.
// Every next time Get is called, it checks the current time against periods,
// depending on the result of the previous produce call.
// If time is expired, produce is called again, otherwise the cached value is returned.
//
// Neither ctx nor produce can be nil.
// None of the periods can be negative.
// A zero period is a corner case and means, that produce will be called every time Get is called.
func OnDemand(ctx context.Context, periods periods.Periods, produce ProduceFunc) Renewable {
	if ctx == nil {
		panic(errors.New("context must not be nil"))
	}

	if periods.Default < 0 {
		panic(fmt.Errorf("default period must be zero or positive: %v", periods.Default))
	}

	if periods.Error < 0 {
		panic(fmt.Errorf("error period must be zero or positive: %v", periods.Error))
	}

	if produce == nil {
		panic(errors.New("produce must not be nil"))
	}

	return &onDemand{
		base:    base{ctx: ctx, produce: produce},
		periods: periods,
		lock:    &sync.RWMutex{},
	}
}

func OnDemandNoCtx(periods periods.Periods, produce ProduceFunc) Renewable {
	return OnDemand(context.Background(), periods, produce)
}

// Returns a new instance of Renewable that will use the soft-hard periods strategy.
// First time Renewable.Get is called, produce is called and the result is cached.
// Every next time Get is called, it checks the current time against periods,
// depending on the result of the previous produce call.
//
// If soft deadline is NOT passed, the cached value is returned.
// Otherwise, if hard deadline is NOT passed, the cached value is returned
// and a goroutine is started to update the result asynchronously.
// If both periods are expired, produce is called again and a caller waits for it to complete.
//
// Neither ctx nor produce can be nil.
// None of the periods can be negative.
// None of the hard periods can be less, than corresponding soft period.
// A zero period is a corner case and means, that produce will be called every time Get is called.
func SoftHard(ctx context.Context, soft periods.Periods, hard periods.Periods, produce ProduceFunc) Renewable {
	if ctx == nil {
		panic(errors.New("context must not be nil"))
	}

	if soft.Default < 0 {
		panic(fmt.Errorf("default soft period must be zero or positive: %v", soft.Default))
	}

	if soft.Error < 0 {
		panic(fmt.Errorf("error soft period must be zero or positive: %v", soft.Error))
	}

	if hard.Default < soft.Default {
		panic(fmt.Errorf(
			"default hard period must be equal or greater than soft: %v < %v",
			hard.Default, soft.Default,
		))
	}

	if hard.Error < soft.Error {
		panic(fmt.Errorf(
			"error hard period must be equal or greater than soft: %v < %v",
			hard.Error, soft.Error,
		))
	}

	if produce == nil {
		panic(errors.New("produce must not be nil"))
	}

	return &softHard{
		base:         base{ctx: ctx, produce: produce},
		soft:         soft,
		hard:         hard,
		lock:         &sync.RWMutex{},
		asyncResults: make(chan result),
	}
}

func SoftHardNoCtx(soft periods.Periods, hard periods.Periods, produce ProduceFunc) Renewable {
	return SoftHard(context.Background(), soft, hard, produce)
}

func GetOrPanic(r Renewable) interface{} {
	if value, err := r.Get(); err != nil {
		panic(err)
	} else {
		return value
	}
}
