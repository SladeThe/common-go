package renewable

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/SladeThe/common-go/renewable/periods"
)

/* TODO type softHardValType uint8

const (
	softHardValTypeInvalid softHardValType = iota
	softHardValTypeHard
	softHardValTypeSoft

	softHardAsyncStateNotRunning uint32 = iota
	softHardAsyncStateRunning
	softHardAsyncStateCommitting
)

func (vt softHardValType) isValid() bool {
	return vt > softHardValTypeInvalid
}*/

var _ Renewable = &softHard{}

type softHard struct {
	producing

	soft periods.Periods
	hard periods.Periods

	lock      *sync.Mutex
	condition *sync.Cond
	container *atomic.Value
	state     uint32
}

/*func (r *softHard) startUpdateAsync() {
	go func() {
		println("BEFORE r.results <- result{val: val, err: err}")
		val, err := r.produce(r.ctx)
		r.results <- result{val: val, err: err, time: time.Now()} // TODO when now is calculated?
		println("AFTER r.results <- result{val: val, err: err}")
	}()
}*/

//func (r *softHard) finishUpdateAsync() bool {
/* TODO select {
case r.result = <-r.results:
	println("r.state = softHardAsyncStateNotRunning")
	r.state = softHardAsyncStateNotRunning
	r.produceTime = time.Now()
	return true
default:
	return false
}*/
//}

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

/*func (r *softHard) get() (interface{}, error, softHardValType) {
	if raw := r.container.Load(); raw != nil {
		res := raw.(result)
		now := time.Now()

		if res.isValidAt(now, r.soft) {
			return res.val, res.err, softHardValTypeSoft
		}

		if res.isValidAt(now, r.hard) {
			select {
			case res = <-r.results:
				r.container.Store(res)
				atomic.StoreUint32(&r.state, softHardAsyncStateNotRunning)
			default:
				if atomic.CompareAndSwapUint32(&r.state, softHardAsyncStateNotRunning, softHardAsyncStateRunning) {
					r.startUpdateAsync()
				}
			}

			return res.val, res.err, softHardValTypeHard
		}
	}

	if atomic.CompareAndSwapUint32(&r.state, softHardAsyncStateNotRunning, softHardAsyncStateRunning) {
		val, err := r.produce(r.ctx)
		r.container.Store(result{val: val, err: err, time: time.Now()})
		atomic.StoreUint32(&r.state, softHardAsyncStateNotRunning)
		return val, err, nil
	}

	// TODO
	return nil, nil, softHardValTypeInvalid
}

func (r *softHard) Get() (interface{}, error) {
	r.lock.RLock()
	val, err, valType := r.get()
	// TODO r.lock.RUnlock()

	if valType.isValid() {
		if valType == softHardValTypeHard {
			println("CHECK atomic.CompareAndSwapUint32(&r.state, softHardAsyncStateNotRunning, softHardAsyncStateRunning)")
			if atomic.CompareAndSwapUint32(&r.state, softHardAsyncStateNotRunning, softHardAsyncStateRunning) {
				println("TRUE atomic.CompareAndSwapUint32(&r.state, softHardAsyncStateNotRunning, softHardAsyncStateRunning)")
				r.startUpdateAsync()
				println("r.lock.RUnlock()")
				r.lock.RUnlock()
			} else {
				println("FALSE atomic.CompareAndSwapUint32(&r.state, softHardAsyncStateNotRunning, softHardAsyncStateRunning)")
				select {
				case result := <-r.results:
					println("GOT result = <-r.results")
					if !atomic.CompareAndSwapUint32(&r.state, softHardAsyncStateRunning, softHardAsyncStateCommitting) {
						panic("!atomic.CompareAndSwapUint32(&r.state, softHardAsyncStateRunning, softHardAsyncStateCommitting)") // TODO
					}
					println("r.lock.RUnlock()")
					r.lock.RUnlock()

					println("r.lock.Lock()")
					r.lock.Lock()
					println("CHECK 1 r.state == softHardAsyncStateCommitting")
					if r.state == softHardAsyncStateCommitting {
						println("TRUE 1 r.state == softHardAsyncStateCommitting")
						r.state = softHardAsyncStateNotRunning
						r.result = result
						r.produceTime = time.Now()
					} else {
						println("FALSE 1 r.state == softHardAsyncStateCommitting")
					}
					val, err = r.val, r.err
					println("r.lock.Unlock()")
					r.lock.Unlock()
				default:
					println("DEFAULT r.lock.RUnlock()")
					r.lock.RUnlock()
				}
			}
		}

		return val, err
	}

	r.lock.RUnlock()
	r.lock.Lock()
	defer r.lock.Unlock()

	for r.state == softHardAsyncStateCommitting {
		// TODO
	}

	if val, err, valType = r.get(); valType.isValid() {
		if valType == softHardValTypeHard {
			println("BEFORE r.updateAsync()")
			if r.state == softHardAsyncStateRunning {
				if r.finishUpdateAsync() {
					val, err = r.val, r.err
				}
			} else {
				r.startUpdateAsync()
			}
			println("AFTER r.updateAsync()")
		}

		return val, err
	}

	if r.asyncRunning {
		println("BEFORE r.result = <-r.results")
		r.result = <-r.results
		r.asyncRunning = false
		println("AFTER r.result = <-r.results")
	} else {
		println("BEFORE r.val, r.err = r.produce(r.ctx)")
		r.val, r.err = r.produce(r.ctx)
		println("AFTER r.val, r.err = r.produce(r.ctx)")
	}

	r.produceTime = time.Now()
	return r.val, r.err
}*/
