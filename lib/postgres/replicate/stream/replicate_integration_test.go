package stream_test

import (
	"net"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/docker"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/database"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream"
	dockerTester "github.com/pagarme/warp-pipe/lib/tester/docker"
)

func TestIntegrationStreamReplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	t.Run("ConnectNetError", func(t *testing.T) {

		var err error

		r := stream.New(postgres.Config{
			Host:     "localhost",
			Port:     postgres.DefaultPort,
			User:     postgres.DefaultUser,
			Password: "password",
			Database: "database",
		})

		err = r.Start()
		require.Error(t, err)
		require.IsType(t, &net.OpError{}, errors.Cause(err))
	})

	// setup postgres server container
	ipAddress, deferFn := dockerTester.Run(t, docker.Config{
		WaitTimeout: docker.DefaultWaitTimeout,
		URL:         "warp-pipe",
		Image:       "postgres-server",
		Tag:         "9.5.6",
	})
	defer deferFn()

	pgConfig := postgres.Config{
		Host:     ipAddress,
		Port:     postgres.DefaultPort,
		User:     postgres.DefaultUser,
		Database: "test-replicate",
		Password: "postgres",

		Slot:   "test_replicate_slot",
		Plugin: "test_decoding",
		Driver: "pgx",

		ConnectTimeout: 10 * time.Second,

		CreateDatabaseIfNotExist: true,
	}

	// setup database
	db := database.New(pgConfig)
	require.NoError(t, db.Connect())
	defer db.Disconnect()

	t.Run("FullStart", func(t *testing.T) {

		var err error

		r := stream.New(pgConfig)
		err = r.Start()
		require.NoError(t, err)
	})
}
