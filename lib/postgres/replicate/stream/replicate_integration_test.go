package stream_test

import (
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream"
	postgresTester "github.com/pagarme/warp-pipe/lib/tester/postgres"
)

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

		started, err := r.Start()
		require.NoError(t, err)
		require.True(t, started)
	})
}
