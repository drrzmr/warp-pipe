package collector_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/pipeline/collector"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

type testCollector struct {
	collector.Collector
	t *testing.T
	n uint64
}

const numberOfMessages = uint64(100)

func newTestCollector(t *testing.T) *testCollector {
	return &testCollector{
		t: t,
		n: numberOfMessages, // number of messages
	}
}

func (c *testCollector) Init() (err error) { return nil }

func (c *testCollector) Collect(publishCh chan<- message.Message) {
	defer close(publishCh)

	for i := uint64(0); i < c.n; i++ {
		publishCh <- message.New(message.Payload{
			"data": i,
		})
	}
}

func (c *testCollector) UpdateOffset(offsetCh <-chan uint64) {

	expectedOffset := uint64(0)

	for offset := range offsetCh {
		require.Equal(c.t, expectedOffset, offset)
		expectedOffset++
	}
}

func TestCollectorRun(t *testing.T) {

	publishCh, offsetCh, err := collector.Run(newTestCollector(t))
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
