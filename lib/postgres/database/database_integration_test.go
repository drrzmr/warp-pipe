package database_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/docker"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/database"
	dockerTester "github.com/pagarme/warp-pipe/lib/tester/docker"
	"github.com/pagarme/warp-pipe/lib/tester/file"
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

		err = d.Disconnect()
		require.NoError(t, err)
	})

	t.Run("CreateDatabase", func(t *testing.T) {

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
			Database: "test",
			Password: "postgres",

			Driver:         "pgx",
			ConnectTimeout: postgres.MinConnectTimeout,

			CreateDatabaseIfNotExist: true,
		})

		err = d.Connect()
		require.NoError(t, err)
		defer func() {
			err = d.Disconnect()
			require.NoError(t, err)
		}()

		db := d.DB()
		require.NotNil(t, db)

		file.Config.InputExtension = "sql"

		_, err = db.Exec(file.Load(t, "test"))
		require.NoError(t, err)

		rows, err := db.Query("SELECT * FROM test")
		require.NoError(t, err)
		defer func() {
			err = rows.Close()
			require.NoError(t, err)
		}()

		cols, err := rows.Columns()
		require.NoError(t, err)
		require.Equal(t, []string{"a", "b", "c"}, cols)

		colsTypes, err := rows.ColumnTypes()
		require.NoError(t, err)
		require.Equal(t, len(cols), len(colsTypes))

		type row struct {
			a int64
			b string
			c time.Time
		}

		var table []row
		for rows.Next() {
			r := row{}
			err = rows.Scan(&r.a, &r.b, &r.c)
			require.NoError(t, err)
			table = append(table, r)
		}

		for i, row := range table {
			expected := int64(i) + 1
			require.Equal(t, expected, row.a)
			require.Equal(t, fmt.Sprintf("test %d", expected), row.b)
		}
	})

}
