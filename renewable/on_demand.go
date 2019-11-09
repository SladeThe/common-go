package renewable

import (
	"sync"
	"time"
)

var _ Renewable = &onDemand{}

type onDemand struct {
	result

	period      time.Duration
	errPeriod   time.Duration
	produce     ProduceFunc
	produceTime time.Time
	lock        *sync.RWMutex
}

func (r *onDemand) get() (interface{}, error, bool) {
	var period time.Duration

	if r.err == nil {
		period = r.period
	} else {
		period = r.errPeriod
	}

	if !r.produceTime.IsZero() && r.produceTime.Add(period).After(time.Now()) {
		return r.value, r.err, true
	}

	return nil, nil, false
}

func (r *onDemand) Get() (interface{}, error) {
	r.lock.RLock()
	value, err, ok := r.get()
	r.lock.RUnlock()

	if ok {
		return value, err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if value, err, ok = r.get(); ok {
		return value, err
	}

	r.value, r.err = r.produce()
	r.produceTime = time.Now()
	return r.value, r.err
}
