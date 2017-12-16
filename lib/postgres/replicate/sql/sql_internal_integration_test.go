package sql

import (
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/docker"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/retry"
	postgresTester "github.com/pagarme/warp-pipe/lib/tester/postgres"
)

var (
	dockerConfig = docker.Config{
		WaitTimeout: docker.DefaultWaitTimeout,
		URL:         "warp-pipe",
		Image:       "postgres-server",
		Tag:         "9.5.6",
	}
	postgresConfig = postgres.Config{
		Host:     "none.host",
		Port:     postgres.DefaultPort,
		User:     postgres.DefaultUser,
		Database: "test-replicate",
		Password: "postgres",
		Replicate: postgres.ReplicateConfig{
			Slot:   "test_replicate_slot",
			Plugin: "test_decoding",
		},
		SQL: postgres.SQLConfig{
			Driver:                   "pgx",
			ConnectTimeout:           10 * time.Second,
			CreateDatabaseIfNotExist: true,
		},
	}
)

func TestIntegrationSQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	normal, deferFn := postgresTester.DockerRun(t, dockerConfig, &postgresConfig)
	defer deferFn()

	var (
		dsn, _  = postgresConfig.DSN(true, true)
		driver  = postgresConfig.SQL.Driver
		timeout = postgresConfig.SQL.ConnectTimeout
		slot    = postgresConfig.Replicate.Slot
		plugin  = postgresConfig.Replicate.Plugin
		db      *sqlx.DB

		normalDB = normal.DB()
	)

	_, err := normalDB.Exec(`
CREATE TABLE test
(
  id        SERIAL,
  name      VARCHAR(30),
  timestamp TIMESTAMP NOT NULL,
  PRIMARY KEY (id)
);`)
	require.NoError(t, err)

	err, innerErr := retry.Do(timeout, func() (err error) {
		db, err = sqlx.Connect(driver, dsn)
		return err
	})
	require.NoError(t, innerErr)
	require.NoError(t, err)

	t.Run("createSlot", func(t *testing.T) {
		created, err := createSlot(db, slot, plugin)
		require.NoError(t, err)
		require.True(t, created)
	})

	t.Run("listSlots", func(t *testing.T) {
		slots, err := listSlots(db)
		require.NoError(t, err)
		require.Len(t, slots, 1)
		require.Equal(t, slot, slots[0].SlotName)
		require.Equal(t, plugin, slots[0].Plugin)
		require.Equal(t, "logical", slots[0].SlotType)
		require.Equal(t, postgresConfig.Database, slots[0].Database)
		require.False(t, slots[0].Active)
		require.Equal(t, int64(-1), slots[0].ActivePID)
		require.NotEmpty(t, slots[0].RestartLSN)
	})

	t.Run("getAllChanges", func(t *testing.T) {
		changes, err := getAllChanges(db, slot)
		require.NoError(t, err)
		require.Len(t, changes, 0)

		_, err = normalDB.Exec("INSERT INTO test (name, timestamp) VALUES ('test1', now());")
		require.NoError(t, err)

		changes, err = getAllChanges(db, slot)
		require.NoError(t, err)
		require.Len(t, changes, 3)

		var (
			begin     = changes[0]
			operation = changes[1]
			commit    = changes[2]
		)

		require.True(t, strings.HasPrefix(begin.Data, "BEGIN "))
		operationPrefix := "table public.test: INSERT: " +
			"id[integer]:1 " +
			"name[character varying]:'test1' " +
			"\"timestamp\"[timestamp without time zone]:'"
		require.True(t, strings.HasPrefix(operation.Data, operationPrefix))
		require.True(t, strings.HasPrefix(commit.Data, "COMMIT "))
		require.True(t, operation.XID == begin.XID && operation.XID == commit.XID)

		changes, err = getAllChanges(db, slot)
		require.NoError(t, err)
		require.Len(t, changes, 0)
	})
}
