package sql

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	postgresTester "github.com/pagarme/warp-pipe/lib/tester/postgres"
)

func TestInternalIntegrationSQLReplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	dockerConfig := replicate.CreateTestDockerConfig(t)
	postgresConfig := replicate.CreateTestPostgresConfig(t)

	_, deferFn := postgresTester.DockerRun(t, dockerConfig, &postgresConfig)
	defer deferFn()

	r := New(postgresConfig)
	require.NotNil(t, r)

	t.Run("connect", func(t *testing.T) {
		connected := r.isConnected()
		require.False(t, connected)

		err := r.connect()
		require.NoError(t, err)

		connected = r.isConnected()
		require.True(t, connected)
	})
}
