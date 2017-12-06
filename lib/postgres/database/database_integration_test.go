package database_test

import (
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/database"
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
}
