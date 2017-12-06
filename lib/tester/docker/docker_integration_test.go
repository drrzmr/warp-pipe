package docker_test

import (
	"testing"

	"github.com/pagarme/warp-pipe/lib/docker"
	dockerTester "github.com/pagarme/warp-pipe/lib/tester/docker"
	"github.com/pagarme/warp-pipe/lib/tester/ping"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDocker(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	ipAddress, deferFn := dockerTester.Run(t, docker.Config{
		WaitTimeout: docker.DefaultWaitTimeout,
		URL:         "warp-pipe",
		Image:       "postgres-server",
		Tag:         "9.5.6",
	})
	defer deferFn()

	tx, rx := ping.Run(t, ipAddress)
	require.Equal(t, ping.Config.Count, tx)
	require.Equal(t, ping.Config.Count, rx)
}
