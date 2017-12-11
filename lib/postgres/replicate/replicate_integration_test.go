package replicate_test

import (
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
)

func TestIntegrationReplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	t.Run("ConnectNetError", func(t *testing.T) {

		var err error

		r := replicate.New(postgres.Config{
			Host:     "localhost",
			Port:     postgres.DefaultPort,
			User:     postgres.DefaultUser,
			Password: "password",
			Database: "database",
		})

		err = r.Start()
		require.Error(t, err)
		require.IsType(t, &net.OpError{}, errors.Cause(err))
	})
}
