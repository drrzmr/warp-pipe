package replicate_test

import (
	"testing"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	"github.com/stretchr/testify/require"
)

func TestReplicate(t *testing.T) {

	t.Run("Config", func(t *testing.T) {

		var err error

		r := replicate.New(postgres.Config{})
		require.NotNil(t, r)

		err = r.Start()
		require.NoError(t, err)

		config := r.Config()
		require.NotNil(t, config)
	})
}
