package stream_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream"
)

func TestStreamReplicate(t *testing.T) {

	t.Run("Config", func(t *testing.T) {

		r := stream.New(postgres.Config{})
		require.NotNil(t, r)

		config := r.Config()
		require.NotNil(t, config)
	})
}
