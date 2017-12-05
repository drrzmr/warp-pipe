package database_test

import (
	"testing"

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
}
