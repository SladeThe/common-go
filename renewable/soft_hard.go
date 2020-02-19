package renewable

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/SladeThe/common-go/renewable/periods"
)

var _ Renewable = (*softHard)(nil)

type softHard struct {
	producing

	soft periods.Periods
	hard periods.Periods

	lock      *sync.Mutex
	condition *sync.Cond
	container *atomic.Value
	state     uint32
}

func (r *softHard) updateAsync() {
	go func() {
		val, err := r.produce(r.ctx)
		r.container.Store(result{val: val, err: err, time: time.Now()})

		r.lock.Lock()
		defer r.lock.Unlock()

		atomic.StoreUint32(&r.state, 0)
		r.condition.Broadcast()
	}()
}

func (r *softHard) Get() (interface{}, error) {
	if raw := r.container.Load(); raw != nil {
		res := raw.(result)
		now := time.Now()

		if res.isValidAt(now, r.soft) {
			return res.val, res.err
		}

		if res.isValidAt(now, r.hard) {
			if atomic.CompareAndSwapUint32(&r.state, 0, 1) {
				if res = r.container.Load().(result); res.isValidAt(now, r.soft) {
					r.lock.Lock()
					defer r.lock.Unlock()

					atomic.StoreUint32(&r.state, 0)
					r.condition.Broadcast()
				} else {
					r.updateAsync()
				}
			}

			return res.val, res.err
		}
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	for !atomic.CompareAndSwapUint32(&r.state, 0, 1) {
		for atomic.LoadUint32(&r.state) > 0 {
			r.condition.Wait()
		}
	}

	if raw := r.container.Load(); raw != nil {
		res := raw.(result)
		now := time.Now()

		if res.isValidAt(now, r.soft) {
			atomic.StoreUint32(&r.state, 0)
			r.condition.Broadcast()
			return res.val, res.err
		}

		if res.isValidAt(now, r.hard) {
			r.updateAsync()
			return res.val, res.err
		}
	}

	val, err := r.produce(r.ctx)
	r.container.Store(result{val: val, err: err, time: time.Now()})

	atomic.StoreUint32(&r.state, 0)
	r.condition.Broadcast()

	return val, err
}
