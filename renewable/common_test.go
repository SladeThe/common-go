package renewable

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	period          = 75 * time.Millisecond
	errPeriod       = 100 * time.Millisecond
	safeCheckPeriod = 25 * time.Millisecond
)

// Various const limitations.
//
//noinspection GoBoolExpressions
var (
	_ = map[bool]struct{}{
		false:               {},
		safeCheckPeriod > 0: {},
	}

	_ = map[bool]struct{}{
		false:                       {},
		period >= safeCheckPeriod*3: {},
	}

	_ = map[bool]struct{}{
		false:                          {},
		errPeriod >= safeCheckPeriod*3: {},
	}

	_ = map[bool]struct{}{
		false: {},
		period > errPeriod && period-errPeriod >= safeCheckPeriod ||
			errPeriod > period && errPeriod-period >= safeCheckPeriod: {},
	}
)

func simpleTestRenewable(t *testing.T, createRenewable func(produce ProduceFunc) Renewable) {
	const iterCount = 10

	results := make([]result, 0, iterCount*2)
	for i := 0; i < iterCount; i++ {
		results = append(results, result{value: i}, result{err: fmt.Errorf("%d", i)})
	}

	renewable := createRenewable(func() (value interface{}, err error) {
		if len(results) <= 0 {
			assert.FailNow(t, "results are unexpectedly empty")
		}

		r := results[0]
		results = results[1:]
		return r.value, r.err
	})

	time.Sleep(safeCheckPeriod)
	assert.Equal(t, iterCount*2, len(results))

	for i := 0; i < iterCount; i++ {
		value, err := renewable.Get()
		assert.Equal(t, i, value)
		assert.Nil(t, err)
		assert.Equal(t, iterCount*2-i*2-1, len(results))

		if i > 0 {
			time.Sleep(period - 2*safeCheckPeriod)
		} else {
			time.Sleep(period - safeCheckPeriod)
		}

		value, err = renewable.Get()
		assert.Equal(t, i, value)
		assert.Nil(t, err)
		assert.Equal(t, iterCount*2-i*2-1, len(results))

		time.Sleep(2 * safeCheckPeriod)

		value, err = renewable.Get()
		assert.Nil(t, value)
		assert.Equal(t, fmt.Errorf("%d", i), err)
		assert.Equal(t, iterCount*2-i*2-2, len(results))

		time.Sleep(errPeriod - 2*safeCheckPeriod)

		value, err = renewable.Get()
		assert.Nil(t, value)
		assert.Equal(t, fmt.Errorf("%d", i), err)
		assert.Equal(t, iterCount*2-i*2-2, len(results))

		time.Sleep(2 * safeCheckPeriod)
	}
}

func asyncTestRenewable(t *testing.T, createRenewable func(produce ProduceFunc) Renewable) {
	const (
		iterCount         = 10
		getRoutineCount   = 2
		checkRoutineCount = 4
	)

	var iter uint64 = 0

	renewable := createRenewable(func() (value interface{}, err error) {
		i := atomic.LoadUint64(&iter)
		if i > (iterCount+1)*2 {
			assert.FailNow(t, fmt.Sprintf("iter is unexpectedly large: %d", i))
		}
		defer atomic.StoreUint64(&iter, i+1)

		if i%2 == 0 {
			return i / 2, nil
		} else {
			return nil, fmt.Errorf("%d", (i-1)/2)
		}
	})

	time.Sleep(safeCheckPeriod)
	assert.Equal(t, uint64(0), iter)

	ctx, cancel := context.WithCancel(context.Background())

	var gwg sync.WaitGroup
	gwg.Add(getRoutineCount)

	for i := 0; i < getRoutineCount; i++ {
		go func() {
			assert.NotPanics(t, func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						_, _ = renewable.Get()
					}
				}
			})

			gwg.Done()
		}()
	}

	var cwg sync.WaitGroup
	cwg.Add(checkRoutineCount)

	for i := 0; i < checkRoutineCount; i++ {
		go func() {
			assert.NotPanics(t, func() {
				for i := 0; i < iterCount; i++ {
					value, err := renewable.Get()
					assert.Equal(t, uint64(i), value)
					assert.Nil(t, err)

					if i > 0 {
						time.Sleep(period - 2*safeCheckPeriod)
					} else {
						time.Sleep(period - safeCheckPeriod)
					}

					value, err = renewable.Get()
					assert.Equal(t, uint64(i), value)
					assert.Nil(t, err)

					time.Sleep(2 * safeCheckPeriod)

					value, err = renewable.Get()
					assert.Nil(t, value)
					assert.Equal(t, fmt.Errorf("%d", i), err)

					time.Sleep(errPeriod - 2*safeCheckPeriod)

					value, err = renewable.Get()
					assert.Nil(t, value)
					assert.Equal(t, fmt.Errorf("%d", i), err)

					if i < iterCount-1 {
						time.Sleep(2 * safeCheckPeriod)
					}
				}
			})

			cwg.Done()
		}()
	}

	cwg.Wait()
	cancel()
	gwg.Wait()
}
