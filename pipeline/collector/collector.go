package collector

import (
	"github.com/pagarme/warp-pipe/pipeline/message"
	"github.com/pkg/errors"
)

// Collector interface
type Collector interface {
	Init() (err error)
	Collect(publishCh chan<- message.Message)
	UpdateOffset(updateOffsetCh <-chan uint64)
}

// Run collector
func Run(c Collector) (publishCh <-chan message.Message, offsetCh chan<- uint64, err error) {

	var (
		publish = make(chan message.Message)
		offset  = make(chan uint64)
	)

	if err = c.Init(); err != nil {
		return nil, nil, errors.Wrap(err, "Could not initialize stage collector")
	}

	go c.Collect(publish)
	go c.UpdateOffset(offset)

	return publish, offset, nil
}
