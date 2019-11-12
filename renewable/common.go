package renewable

import (
	"context"
	"time"
)

type result struct {
	val interface{}
	err error
}

type base struct {
	result

	ctx         context.Context
	produce     ProduceFunc
	produceTime time.Time
}
