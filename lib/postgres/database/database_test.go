package database_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/database"
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
			Driver: "",
		})

		err = d.Connect()
		require.Error(t, err)

		causeErr := errors.Cause(err)
		require.Error(t, causeErr)

		expectedErr := fmt.Errorf("sql: unknown driver %q (forgotten import?)", "")
		require.Equal(t, expectedErr, causeErr)
	})
}
