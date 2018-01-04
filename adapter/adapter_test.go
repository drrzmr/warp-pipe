package adapter_test

import (
	"context"
	"fmt"
	"testing"

	mockCollector "github.com/pagarme/warp-pipe/adapter/collector/mock"
	mockProcessor "github.com/pagarme/warp-pipe/adapter/processor/mock"
	"github.com/pagarme/warp-pipe/pipeline"
	"github.com/stretchr/testify/require"
)

func TestMockAdapter(t *testing.T) {

	collectFunc := func(messageID uint64, publishCh chan<- pipeline.Message) (end bool) {
		publishCh <- pipeline.NewMessage(pipeline.Payload{
			"rawData": messageID,
		})
		fmt.Println("[send message] messageID:", messageID)
		return false
	}

	updateOffsetFunc := func(offset uint64) {
		fmt.Println("[update offset for message] offset:", offset)
	}

	processFunc := func(msg pipeline.Message, out chan<- pipeline.Message) (end bool) {

		var (
			rawData       = msg.Get("rawData").(uint64)
			processedData = rawData * 3
		)
		out <- pipeline.NewMessage(pipeline.Payload{
			"rawData":       rawData,
			"processedData": processedData,
		})

		return false
	}

	var (
		run       = pipeline.NewRunner("adapter_test")
		ctx       = context.Background()
		collector = mockCollector.New(10, collectFunc, updateOffsetFunc)
		processor = mockProcessor.New(processFunc)
	)

	outOfCollector, offset, err := run.Collector(ctx, collector)
	close(offset)
	require.NoError(t, err)

	outOfProcessor, err := run.Processor(ctx, processor, outOfCollector)
	require.NoError(t, err)

	for msg := range outOfProcessor {
		var (
			rawData       = msg.Get("rawData").(uint64)
			processedData = msg.Get("processedData").(uint64)
		)

		fmt.Println("[message] rawData:", rawData, ", processedData", processedData)

		require.Equal(t, rawData*3, processedData)
	}
}
