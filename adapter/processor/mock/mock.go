package mock

import (
	"context"

	"github.com/pagarme/warp-pipe/pipeline"
)

type processFunc func(msg pipeline.Message, outCh chan<- pipeline.Message) (end bool)

// Processor object
type Processor struct {
	pipeline.Processor
	processCb processFunc

	inCh <-chan pipeline.Message
}

// New return a Collector instance
func New(processCb processFunc) *Processor {

	return &Processor{
		processCb: processCb,
	}
}

// Init implements method from interface
func (p *Processor) Init(ctx context.Context, inCh <-chan pipeline.Message) (err error) {

	p.inCh = inCh

	return nil
}

// Process implements method from interface
func (p *Processor) Process(outCh chan<- pipeline.Message) {
	defer close(outCh)

	for msg := range p.inCh {
		p.processCb(msg, outCh)
	}
}
