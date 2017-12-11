package replicate_test

import (
	"testing"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	"github.com/stretchr/testify/require"
)

func TestReplicate(t *testing.T) {

	t.Run("Config", func(t *testing.T) {

		r := replicate.New(postgres.Config{})
		require.NotNil(t, r)
		config := r.Config()
		require.NotNil(t, config)
	})
}
