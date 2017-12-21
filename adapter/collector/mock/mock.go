package mock

import (
	"github.com/pagarme/warp-pipe/pipeline/collector"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

type collectFunc func(messageID uint64, publishCh chan<- message.Message) (end bool)

type updateOffsetFunc func(offset uint64)

// Collector object
type Collector struct {
	collector.Collector
	numOfMessages  uint64
	collectCb      collectFunc
	updateOffsetCb updateOffsetFunc
}

// New return a Collector instance
func New(numberOfMessages uint64, collectCb collectFunc, updateOffsetCb updateOffsetFunc) *Collector {

	return &Collector{
		numOfMessages:  numberOfMessages,
		collectCb:      collectCb,
		updateOffsetCb: updateOffsetCb,
	}
}

// Init implements method from interface
func (c *Collector) Init() (err error) { return nil }

// Collect implements method from interface
func (c *Collector) Collect(publishCh chan<- message.Message) {
	defer close(publishCh)

	for i := uint64(0); i < c.numOfMessages; i++ {
		var end bool
		if end = c.collectCb(i, publishCh); end {
			return
		}
	}
}

// UpdateOffset implements method from interface
func (c *Collector) UpdateOffset(offsetCh <-chan uint64) {
	for offset := range offsetCh {
		c.updateOffsetCb(offset)
	}
}
