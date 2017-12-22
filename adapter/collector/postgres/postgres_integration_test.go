package postgres_test

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/pagarme/warp-pipe/adapter/collector/postgres"
	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream"
	tester "github.com/pagarme/warp-pipe/lib/tester/postgres"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

var logger = log.Development("test")

func TestIntegrationPostgresAdapter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	var (
		dockerConfig   = replicate.CreateTestDockerConfig(t)
		postgresConfig = replicate.CreateTestPostgresConfig(t)
		_, deferFn     = tester.DockerRun(t, dockerConfig, &postgresConfig)
		collector      = postgres.New(stream.New(postgresConfig))
		publishCh      = make(chan message.Message)
		offsetCh       = make(chan uint64)
		done           = make(chan struct{})
		ctx, cancel    = context.WithCancel(context.Background())
	)
	defer deferFn()

	time.AfterFunc(10*time.Second, func() {
		logger.Debug("canceling...")
		cancel()
	})

	collector.Init(ctx)
	go collector.Collect(publishCh)
	go collector.UpdateOffset(offsetCh)
	go func(publishCh <-chan message.Message, offsetCh chan<- uint64) {
		defer close(offsetCh)
		defer close(done)

		for msg := range publishCh {
			logger.Debug("new message",
				zap.Time("timestamp", msg.Timestamp()),
			)
			offsetCh <- 0
		}
	}(publishCh, offsetCh)

	<-done
}
