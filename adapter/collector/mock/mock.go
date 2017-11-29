package mock

import (
	"github.com/pagarme/warp-pipe/pipeline/collector"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

type collectFunc func(c *Collector, messageID uint64, publishCh chan<- message.Message) (end bool)

type updateOffsetFunc func(c *Collector, offset uint64)

type channels struct {
	publish chan<- message.Message
	offset  <-chan uint64
}

// Collector object
type Collector struct {
	collector.Collector
	numOfMessages  uint64
	collectCb      collectFunc
	updateOffsetCb updateOffsetFunc
	channels       *channels
}

// New return a Collector instance
func New(numberOfMessages uint64, collectCb collectFunc, updateOffsetCb updateOffsetFunc) *Collector {

	return &Collector{
		numOfMessages:  numberOfMessages,
		collectCb:      collectCb,
		updateOffsetCb: updateOffsetCb,
		channels:       nil,
	}
}

// Init implements method from interface
func (c *Collector) Init(publishCh chan<- message.Message, offsetCh <-chan uint64) (err error) {

	c.channels = &channels{
		publish: publishCh,
		offset:  offsetCh,
	}
	return nil
}

// Collect implements method from interface
func (c *Collector) Collect() {
	if c.channels == nil {
		panic("Collector not initialized!")
	}

	defer close(c.channels.publish)

	for i := uint64(0); i < c.numOfMessages; i++ {
		var end bool
		if end = c.collectCb(c, i, c.channels.publish); end {
			return
		}
	}
}

// UpdateOffset implements method from interface
func (c *Collector) UpdateOffset() {
	if c.channels == nil {
		panic("Collector not initialized!")
	}

	for offset := range c.channels.offset {
		c.updateOffsetCb(c, offset)
	}
}
