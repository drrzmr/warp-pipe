package mock_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pagarme/warp-pipe/adapter/processor/mock"
	"github.com/pagarme/warp-pipe/pipeline"
	"github.com/stretchr/testify/require"
)

func TestMockProcessor(t *testing.T) {

	process := func(msg pipeline.Message, outCh chan<- pipeline.Message) (end bool) {

		rawData := msg.Get("rawData").(uint64)
		processedData := rawData * 10
		fmt.Printf("[process callback] rawData: %d, processedData: %d\n", rawData, processedData)

		outCh <- pipeline.NewMessage(pipeline.Payload{
			"rawData":       rawData,
			"processedData": processedData,
		})

		return false
	}

	var (
		processor = mock.New(process)
		inCh      = make(chan pipeline.Message)
		outCh     = make(chan pipeline.Message)
		ctx       = context.Background()
	)

	err := processor.Init(ctx, inCh)
	require.NoError(t, err)

	go processor.Process(outCh)

	go func(out chan<- pipeline.Message, n uint64) {
		defer close(out)

		for i := uint64(0); i < n; i++ {
			out <- pipeline.NewMessage(pipeline.Payload{
				"rawData": i,
			})
		}
	}(inCh, 100)

	for msg := range outCh {

		rawData := msg.Get("rawData").(uint64)
		processedData := msg.Get("processedData").(uint64)

		require.Equal(t, rawData*10, processedData)
	}
}
