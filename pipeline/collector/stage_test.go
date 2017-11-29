package collector_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/pipeline/collector"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

type testCollector struct {
	collector.Collector
	t              *testing.T
	n              uint64
	publishCh      chan<- message.Message
	updateOffsetCh <-chan uint64
}

const numberOfMessages = uint64(100)

func newTestCollector(t *testing.T) *testCollector {
	return &testCollector{
		t: t,
		n: numberOfMessages, // number of messages
	}
}

func (c *testCollector) Init(publishCh chan<- message.Message, updateOffsetCh <-chan uint64) (err error) {

	c.publishCh = publishCh
	c.updateOffsetCh = updateOffsetCh

	return nil
}

func (c *testCollector) Collect() {
	defer close(c.publishCh)

	for i := uint64(0); i < c.n; i++ {
		c.publishCh <- message.New(message.Payload{
			"data": i,
		})
	}
}

func (c *testCollector) UpdateOffset() {

	expectedOffset := uint64(0)

	for offset := range c.updateOffsetCh {
		require.Equal(c.t, expectedOffset, offset)
		expectedOffset++
	}
}

func TestStage_Run(t *testing.T) {

	stage := collector.NewStage(newTestCollector(t))
	require.NotNil(t, stage)

	publishCh, offsetCh, err := stage.Run()
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
