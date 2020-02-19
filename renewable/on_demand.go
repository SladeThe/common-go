package renewable

import (
	"sync"
	"time"

	"github.com/SladeThe/common-go/renewable/periods"
)

var _ Renewable = (*onDemand)(nil)

type onDemand struct {
	result
	producing

	periods periods.Periods
	lock    *sync.RWMutex
}

func (r *onDemand) get() (interface{}, error, bool) {
	if r.isValidNow(r.periods) {
		return r.val, r.err, true
	}

	return nil, nil, false
}

func (r *onDemand) Get() (interface{}, error) {
	r.lock.RLock()
	val, err, ok := r.get()
	r.lock.RUnlock()

	if ok {
		return val, err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if val, err, ok = r.get(); ok {
		return val, err
	}

	r.val, r.err = r.produce(r.ctx)
	r.time = time.Now()
	return r.val, r.err
}
