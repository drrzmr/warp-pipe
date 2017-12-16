package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/docker"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/database"
	"github.com/pagarme/warp-pipe/lib/tester"
	dockerTester "github.com/pagarme/warp-pipe/lib/tester/docker"
)

// Connect postgres test helper
func Connect(t *testing.T, config postgres.Config) (db *database.Database, deferFn tester.DeferFunc) {
	t.Helper()

	db = database.New(config)
	err := db.Connect()
	require.NoError(t, err)

	deferFn = func() {
		err = db.Disconnect()
		require.NoError(t, err)
	}

	return db, deferFn
}

// DockerRun postgres test helper
func DockerRun(t *testing.T, dockerConfig docker.Config, postgresConfig *postgres.Config) (
	db *database.Database, deferFn tester.DeferFunc) {

	t.Helper()

	ipAddress, dockerDeferFn := dockerTester.Run(t, dockerConfig)
	postgresConfig.Host = ipAddress
	db, postgresDeferFn := Connect(t, *postgresConfig)

	deferFn = func() {
		postgresDeferFn()
		dockerDeferFn()
	}

	return db, deferFn
}
