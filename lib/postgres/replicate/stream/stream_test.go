package stream_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream"
)

func TestStreamReplicate(t *testing.T) {

	t.Run("Config", func(t *testing.T) {

		s := stream.New(postgres.Config{})
		require.NotNil(t, s)

		config := s.Config()
		require.NotNil(t, config)
	})
}
