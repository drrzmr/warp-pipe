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

// Processor interface
type Processor interface {
	Init(ctx context.Context, inCh <-chan Message) (err error)
	Process(outCh chan<- Message)
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

// Processor runner
func (r Runner) Processor(ctx context.Context, p Processor, inCh <-chan Message) (outCh <-chan Message, err error) {

	ch := make(chan Message)

	if err = p.Init(ctx, inCh); err != nil {
		return nil, errors.Wrap(err, "Could not initialize stage processor")
	}

	go p.Process(ch)

	return ch, nil
}
