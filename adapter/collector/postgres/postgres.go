package postgres

import (
	"go.uber.org/zap"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream"
	"github.com/pagarme/warp-pipe/pipeline/collector"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

// Collector object
type Collector struct {
	collector.Collector
	stream *stream.Replicate
}

var logger = log.Development("collector")

// New create a new collector
func New(stream *stream.Replicate) *Collector {

	return &Collector{
		stream: stream,
	}
}

// Init method
func (c *Collector) Init() (err error) {

	logger.Debug("--> Init()")
	defer logger.Debug("<-- Init()")

	return nil
}

// Collect method
func (c *Collector) Collect(publishCh chan<- message.Message) {

	logger.Debug("--> Collect()")
	defer logger.Debug("<-- Collect()")

	defer close(publishCh)
}

// UpdateOffset method
func (c *Collector) UpdateOffset(offsetCh <-chan uint64) {

	logger.Debug("--> UpdateOffset()")
	defer logger.Debug("<-- UpdateOffset()")

	for offset := range offsetCh {
		logger.Debug("update offset", zap.Uint64("offset", offset))
	}
}
