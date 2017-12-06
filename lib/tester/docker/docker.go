package docker

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/docker"
	"github.com/pagarme/warp-pipe/lib/tester"
)

// Run run a docker container for given config
func Run(t *testing.T, config docker.Config) (ipAddress string, deferFn tester.DeferFunc) {
	t.Helper()
	var err error

	runner := docker.NewRunner(config)
	require.NotNil(t, runner)

	err = runner.Start()
	require.NoError(t, err)

	ipAddress = runner.IPAddress()
	deferFn = func() {
		err = runner.Stop()
		require.NoError(t, err)
	}

	return ipAddress, deferFn
}
