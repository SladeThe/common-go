package renewable

import (
	"sync"
	"time"

	"github.com/SladeThe/common-go/renewable/periods"
)

type softHardValType uint8

const (
	softHardValTypeInvalid softHardValType = iota
	softHardValTypeHard
	softHardValTypeSoft
)

func (vt softHardValType) isValid() bool {
	return vt > softHardValTypeInvalid
}

var _ Renewable = &softHard{}

type softHard struct {
	base

	soft periods.Periods
	hard periods.Periods
	lock *sync.RWMutex

	asyncResults chan result
	asyncRunning bool
}

func (r *softHard) updateAsync() {
	if !r.asyncRunning {
		r.asyncRunning = true

		go func() {
			val, err := r.produce(r.ctx)
			r.asyncResults <- result{val: val, err: err}
		}()
	}
}

func (r *softHard) get() (interface{}, error, softHardValType) {
	if !r.produceTime.IsZero() {
		now := time.Now()

		if r.produceTime.Add(r.soft.Period(r.err)).After(now) {
			return r.val, r.err, softHardValTypeSoft
		}

		if r.produceTime.Add(r.hard.Period(r.err)).After(now) {
			return r.val, r.err, softHardValTypeHard
		}
	}

	return nil, nil, softHardValTypeInvalid
}

func (r *softHard) Get() (interface{}, error) {
	r.lock.RLock()
	val, err, valType := r.get()
	r.lock.RUnlock()

	if valType.isValid() {
		if valType == softHardValTypeHard {
			r.lock.Lock()
			r.updateAsync()
			r.lock.Unlock()
		}

		return val, err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if val, err, valType = r.get(); valType.isValid() {
		if valType == softHardValTypeHard {
			r.updateAsync()
		}

		return val, err
	}

	if r.asyncRunning {
		r.result = <-r.asyncResults
		r.asyncRunning = false
	} else {
		r.val, r.err = r.produce(r.ctx)
	}

	r.produceTime = time.Now()
	return r.val, r.err
}
