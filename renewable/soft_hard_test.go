package renewable

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
	}, asyncTestSoftHardRenewableOnce)
}

func asyncTestSoftHardRenewableOnce(t *testing.T, createRenewable func(produce ProduceFunc) Renewable) {
	const (
		iterCount              = 10
		busyGetRoutineCount    = 2
		valueCheckRoutineCount = 4
	)

	var iter uint64 = 0

	renewable := createRenewable(func(context.Context) (value interface{}, err error) {
		i := atomic.LoadUint64(&iter)
		if i > (iterCount+1)*2 {
			assert.FailNow(t, fmt.Sprintf("iter is unexpectedly large: %d", i))
		}
		defer func() {
			if !atomic.CompareAndSwapUint64(&iter, i, i+1) {
				assert.FailNow(t, "iter has been modified concurrently")
			}
		}()

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
	gwg.Add(busyGetRoutineCount)

	for i := 0; i < busyGetRoutineCount; i++ {
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
	cwg.Add(valueCheckRoutineCount)

	for i := 0; i < valueCheckRoutineCount; i++ {
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

// TODO add more tests
