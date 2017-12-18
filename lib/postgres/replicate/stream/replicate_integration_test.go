package stream_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream"
	postgresTester "github.com/pagarme/warp-pipe/lib/tester/postgres"
)

var logger = log.Development("test")

func TestIntegrationStreamReplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	dockerConfig := replicate.CreateTestDockerConfig(t)
	postgresConfig := replicate.CreateTestPostgresConfig(t)

	t.Run("ConnectNetError", func(t *testing.T) {

		var err error

		r := stream.New(postgresConfig)
		require.NotNil(t, r)

		err = r.Connect()
		require.Error(t, err)
		require.IsType(t, &net.OpError{}, errors.Cause(err))
	})

	// setup postgres server container
	_, deferFn := postgresTester.DockerRun(t, dockerConfig, &postgresConfig)
	defer deferFn()

	t.Run("FullStart", func(t *testing.T) {

		var err error

		r := stream.New(postgresConfig)
		require.NotNil(t, r)

		err = r.Connect()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(10*time.Second, func() {
			logger.Debug("canceling...")
			defer logger.Debug("cancel done")

			cancel()
		})

		started, err := r.Start(ctx, stream.NewDefaultEventListener(r, stream.MockEventHandler))
		require.Error(t, err)
		require.Equal(t, context.Canceled, errors.Cause(err))
		require.False(t, started)
	})
}
