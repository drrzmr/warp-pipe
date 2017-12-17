package sql_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/sql"
	postgresTester "github.com/pagarme/warp-pipe/lib/tester/postgres"
)

func TestIntegrationReplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	postgresConfig := replicate.CreateTestPostgresConfig(t)
	dockerConfig := replicate.CreateTestDockerConfig(t)

	_, deferFn := postgresTester.DockerRun(t, dockerConfig, &postgresConfig)
	defer deferFn()

	r := sql.New(postgresConfig)
	require.NotNil(t, r)

	err := r.Connect()
	require.NoError(t, err)

	created, err := r.CreateSlot()
	require.NoError(t, err)
	require.True(t, created)

	infoList, err := r.ListSlots()
	require.NoError(t, err)
	require.NotNil(t, infoList)
	require.Len(t, infoList, 1)
	require.Equal(t, postgresConfig.Replicate.Slot, infoList[0].SlotName)

	changes, err := r.GetAllChanges()
	require.NoError(t, err)
	require.Len(t, changes, 0)

	err = r.DropSlot()
	require.NoError(t, err)

	err = r.Close()
	require.NoError(t, err)
}
