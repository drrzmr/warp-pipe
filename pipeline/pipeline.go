package pipeline

import (
	"context"

	"github.com/pkg/errors"
)

// Collector interface
type Collector interface {
	Init(ctx context.Context) (err error)
	Collect(publishCh chan<- Message)
	UpdateOffset(updateOffsetCh <-chan uint64)
}

// Runner pipeline stages
type Runner func() (name string)

// NewRunner create a new runner func
func NewRunner(name string) Runner {
	return Runner(func() string {
		return name
	})
}

// Collector runner
func (r Runner) Collector(ctx context.Context, c Collector) (outCh <-chan Message, offsetCh chan<- uint64, err error) {

	var (
		publish = make(chan Message)
		offset  = make(chan uint64)
	)

	if err = c.Init(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "Could not initialize stage collector")
	}

	go c.Collect(publish)
	go c.UpdateOffset(offset)

	return publish, offset, nil
}
