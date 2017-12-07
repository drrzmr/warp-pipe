package retry_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/retry"
)

func TestDo(t *testing.T) {

	t.Run("ZeroTimeout", func(t *testing.T) {

		count := 0
		err, innerErr := retry.Do(0, func() (err error) {

			count++
			return retry.ErrContinue
		})

		require.Error(t, err)
		require.Equal(t, retry.ErrTimeout, err)
		require.Error(t, innerErr)
		require.Equal(t, retry.ErrTimeout, err)
		require.Equal(t, 0, count)
	})
}
