package database_test

import (
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/docker"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/database"
	dockerTester "github.com/pagarme/warp-pipe/lib/tester/docker"
)

func TestDatabaseIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	t.Run("ConnectNetError", func(t *testing.T) {

		var err error

		d := database.New(postgres.Config{
			Host:     "localhost",
			Port:     postgres.DefaultPort,
			User:     postgres.DefaultUser,
			Password: "password",
			Database: "database",

			Driver:         "pgx",
			ConnectTimeout: postgres.MinConnectTimeout,
		})

		err = d.Connect()
		require.Error(t, err)
		require.IsType(t, &net.OpError{}, errors.Cause(err))
	})

	t.Run("ConnectSuccessful", func(t *testing.T) {

		ipAddress, deferFn := dockerTester.Run(t, docker.Config{
			WaitTimeout: docker.DefaultWaitTimeout,
			URL:         "warp-pipe",
			Image:       "postgres-server",
			Tag:         "9.5.6",
		})
		defer deferFn()

		var err error

		d := database.New(postgres.Config{
			Host:     ipAddress,
			Port:     postgres.DefaultPort,
			User:     postgres.DefaultUser,
			Database: "postgres",
			Password: "postgres",

			Driver:         "pgx",
			ConnectTimeout: postgres.MinConnectTimeout,
		})

		err = d.Connect()
		require.NoError(t, err)
	})
}
