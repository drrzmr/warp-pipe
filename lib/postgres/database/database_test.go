package database_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/database"
	"github.com/pagarme/warp-pipe/lib/retry"
)

func TestDatabase(t *testing.T) {

	t.Run("Config", func(t *testing.T) {

		d := database.New(postgres.Config{})
		require.NotNil(t, d)
		config := d.Config()
		require.NotNil(t, config)
	})

	t.Run("ConnectWithEmptyDriver", func(t *testing.T) {

		var err error

		d := database.New(postgres.Config{
			SQL: postgres.SQLConfig{
				Driver: "",
			},
		})

		err = d.Connect()
		require.Error(t, err)

		causeErr := errors.Cause(err)
		require.Error(t, causeErr)

		require.Equal(t, postgres.ErrInvalidDSN, causeErr)
	})

	t.Run("ConnectTimeout", func(t *testing.T) {

		var err error

		d := database.New(postgres.Config{
			Host:     "host",
			Database: "db",
			User:     "user",
			Port:     123,
			Password: "password",
			SQL: postgres.SQLConfig{
				Driver:         "pgx",
				ConnectTimeout: 0,
			},
		})

		err = d.Connect()
		require.Error(t, err)
		require.Equal(t, retry.ErrTimeout, errors.Cause(err))
	})
}
