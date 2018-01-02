package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/pipeline"
)

/*************
 * Collector *
 *************/
type collector struct {
	pipeline.Collector
	t *testing.T
}

const numberOfMessages = uint64(100)

func (c *collector) Init(ctx context.Context) (err error) { return nil }

func (c *collector) Collect(publishCh chan<- pipeline.Message) {
	defer close(publishCh)

	for i := uint64(0); i < numberOfMessages; i++ {
		publishCh <- pipeline.NewMessage(pipeline.Payload{
			"data": i,
		})
	}
}

func (c *collector) UpdateOffset(offsetCh <-chan uint64) {

	expectedOffset := uint64(0)

	for offset := range offsetCh {
		require.Equal(c.t, expectedOffset, offset)
		expectedOffset++
	}
}

func TestCollector(t *testing.T) {

	run := pipeline.NewRunner("test")
	ctx := context.Background()
	publishCh, offsetCh, err := run.Collector(ctx, &collector{t: t})
	require.NoError(t, err)

	expectedMessageData := uint64(0)
	for msg := range publishCh {
		data := msg.Get("data").(uint64)
		require.Equal(t, expectedMessageData, data)
		expectedMessageData++

		offsetCh <- data
	}

	require.Equal(t, numberOfMessages, expectedMessageData)
}
