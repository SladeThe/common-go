package renewable

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	defaultPeriod   = 75 * time.Millisecond
	errorPeriod     = 100 * time.Millisecond
	safeCheckPeriod = 25 * time.Millisecond
)

// Various const limitations.
//
//noinspection GoBoolExpressions
var (
	_ = map[bool]struct{}{false: {}, safeCheckPeriod > 0: {}}
	_ = map[bool]struct{}{false: {}, defaultPeriod >= safeCheckPeriod*3: {}}
	_ = map[bool]struct{}{false: {}, errorPeriod >= safeCheckPeriod*3: {}}

	_ = map[bool]struct{}{
		defaultPeriod > errorPeriod && defaultPeriod-errorPeriod >= safeCheckPeriod: {},
		errorPeriod > defaultPeriod && errorPeriod-defaultPeriod >= safeCheckPeriod: {},
	}
)

func simpleTestRenewable(t *testing.T, createRenewable func(produce ProduceFunc) Renewable) {
	const iterCount = 10

	results := make([]result, 0, iterCount*2)
	for i := 0; i < iterCount; i++ {
		results = append(results, result{val: i}, result{err: fmt.Errorf("%d", i)})
	}

	renewable := createRenewable(func(context.Context) (value interface{}, err error) {
		if len(results) <= 0 {
			assert.FailNow(t, "results are unexpectedly empty")
		}

		r := results[0]
		results = results[1:]
		return r.val, r.err
	})

	time.Sleep(safeCheckPeriod)
	assert.Equal(t, iterCount*2, len(results))

	for i := 0; i < iterCount; i++ {
		value, err := renewable.Get()
		assert.Equal(t, i, value)
		assert.Nil(t, err)
		assert.Equal(t, iterCount*2-i*2-1, len(results))

		if i > 0 {
			time.Sleep(defaultPeriod - 2*safeCheckPeriod)
		} else {
			time.Sleep(defaultPeriod - safeCheckPeriod)
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

		time.Sleep(errorPeriod - 2*safeCheckPeriod)

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
		minNumCPU         = 4
	)

	assert.GreaterOrEqual(t, runtime.NumCPU(), minNumCPU, "insufficient CPUs to run the test properly")
	if t.Failed() {
		t.FailNow()
	}

	var iter uint64 = 0

	renewable := createRenewable(func(context.Context) (value interface{}, err error) {
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
			defer gwg.Done()

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
		}()
	}

	var cwg sync.WaitGroup
	cwg.Add(checkRoutineCount)

	for i := 0; i < checkRoutineCount; i++ {
		go func() {
			defer cwg.Done()

			assert.NotPanics(t, func() {
				for i := 0; i < iterCount; i++ {
					value, err := renewable.Get()
					assert.Equal(t, uint64(i), value)
					assert.Nil(t, err)

					if i > 0 {
						time.Sleep(defaultPeriod - 2*safeCheckPeriod)
					} else {
						time.Sleep(defaultPeriod - safeCheckPeriod)
					}

					value, err = renewable.Get()
					assert.Equal(t, uint64(i), value)
					assert.Nil(t, err)

					time.Sleep(2 * safeCheckPeriod)

					value, err = renewable.Get()
					assert.Nil(t, value)
					assert.Equal(t, fmt.Errorf("%d", i), err)

					time.Sleep(errorPeriod - 2*safeCheckPeriod)

					value, err = renewable.Get()
					assert.Nil(t, value)
					assert.Equal(t, fmt.Errorf("%d", i), err)

					if i < iterCount-1 {
						time.Sleep(2 * safeCheckPeriod)
					}
				}
			})
		}()
	}

	cwg.Wait()
	cancel()
	gwg.Wait()
}
