package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres"
)

func TestConfig(t *testing.T) {

	config := postgres.Config{
		Host:     "host",
		Port:     postgres.DefaultPort,
		User:     postgres.DefaultUser,
		Database: "db",

		Driver:         "pgx",
		ConnectTimeout: postgres.MinConnectTimeout,
	}

	dsnNoDatabase, missing := config.DSN(false, true)
	require.Equal(t, "user=postgres host=host port=5432", dsnNoDatabase)
	require.Len(t, missing, 2)
	require.Contains(t, missing, "database")
	require.Contains(t, missing, "password")

	dsnFull, missing := config.DSN(true, true)
	require.Equal(t, "user=postgres host=host port=5432 database=db", dsnFull)
	require.Len(t, missing, 1)
	require.Contains(t, missing, "password")
}
