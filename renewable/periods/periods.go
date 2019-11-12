package periods

import "time"

type Periods struct {
	Default time.Duration
	Error   time.Duration
}

func (p Periods) Period(err error) time.Duration {
	if err == nil {
		return p.Default
	} else {
		return p.Error
	}
}

func Same(period time.Duration) Periods {
	return Periods{Default: period, Error: period}
}
