package postgres

import (
	"context"
	"time"

	"github.com/pagarme/warp-pipe/adapter/collector/postgres/handler"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream"
	"github.com/pagarme/warp-pipe/pipeline/collector"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

// Collector object
type Collector struct {
	collector.Collector
	stream *stream.Stream
	ctx    context.Context
}

var logger *zap.Logger

func init() { log.Register(&logger, "adapter.collector.postgres") }

// New create a new collector
func New(postgresConfig postgres.Config) *Collector {

	return &Collector{
		stream: stream.New(postgresConfig),
	}
}

// Init method
func (c *Collector) Init(ctx context.Context) (err error) {

	logger.Debug("--> Init()")
	defer logger.Debug("<-- Init()")

	if c.isInitialized() {
		logger.Warn("already initialized")
		return nil
	}

	if err = c.stream.Connect(); err != nil {
		logger.Error("stream connect error", zap.Error(err))
		return errors.WithStack(err)
	}

	if err = c.stream.Start(); err != nil {
		logger.Error("stream start error", zap.Error(err))
		return errors.WithStack(err)
	}

	c.ctx = ctx
	return nil
}

// Collect method
func (c *Collector) Collect(publishCh chan<- message.Message) {

	logger.Debug("--> Collect()")
	defer logger.Debug("<-- Collect()")

	defer close(publishCh)

	var (
		handler  = handler.New(publishCh)
		listener = c.stream.NewDefaultEventListener(handler)
	)

	err := listener.Listen(c.ctx)
	logger.Info("<-- Collect()", zap.Error(err))
}

// UpdateOffset method
func (c *Collector) UpdateOffset(offsetCh <-chan uint64) {

	logger.Debug("--> UpdateOffset()")
	defer logger.Debug("<-- UpdateOffset()")

	var (
		period = c.stream.Config().Streaming.SendStandByStatusPeriod
		ticker = time.NewTicker(period)
		offset uint64
	)
	defer ticker.Stop()

	for {
		select {
		case newOffset, ok := <-offsetCh:
			if !ok {
				logger.Info("offsetCh was closed, exiting", zap.Uint64("lastOffset", offset))
				return
			}

			if newOffset > offset {
				offset = newOffset
			}

		case <-c.ctx.Done():
			logger.Info("canceled, exiting...", zap.Error(c.ctx.Err()))
			return

		case <-ticker.C:
			c.stream.SendStandByStatus(offset)
		}
	}
}

func (c *Collector) isInitialized() (initialized bool) { return c.ctx != nil }
