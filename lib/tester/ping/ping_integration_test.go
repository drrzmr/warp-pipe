package ping_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/tester/ping"
)

func TestIntegrationPing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	tx, rx := ping.Run(t, "localhost")
	require.Equal(t, tx, rx)
}
