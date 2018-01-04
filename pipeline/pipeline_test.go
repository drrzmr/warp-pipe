package pipeline_test

import (
	"context"
	"fmt"
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

/*************
 * Processor *
 *************/
type processor struct {
	pipeline.Processor
	t    *testing.T
	inCh <-chan pipeline.Message
}

func (p *processor) Init(ctx context.Context, inCh <-chan pipeline.Message) (err error) {
	p.inCh = inCh
	return nil
}

func (p *processor) Process(outCh chan<- pipeline.Message) {
	defer close(outCh)

	for msg := range p.inCh {

		var (
			rawData       = msg.Get("rawData").(uint64)
			processedData = rawData * 10
		)

		fmt.Printf("[processor] rawData: %d, processedData: %d\n", rawData, processedData)

		outCh <- pipeline.NewMessage(pipeline.Payload{
			"processedData": processedData,
		})
	}
}

func TestProcessor(t *testing.T) {

	mockInCh := make(chan pipeline.Message)
	go func(outCh chan<- pipeline.Message) {
		defer close(outCh)
		for i := uint64(0); i < numberOfMessages; i++ {
			fmt.Println("[mock producer] rawData:", i)
			outCh <- pipeline.NewMessage(pipeline.Payload{
				"rawData": i,
			})
		}
	}(mockInCh)

	run := pipeline.NewRunner("test")
	outCh, err := run.Processor(context.Background(), &processor{t: t}, mockInCh)
	require.NoError(t, err)

	expectedMessageData := uint64(0)
	for msg := range outCh {
		processedData := msg.Get("processedData").(uint64)
		fmt.Println("[mock consumer] processedData:", processedData)

		require.Equal(t, expectedMessageData*10, processedData)
		expectedMessageData++
	}

	require.Equal(t, numberOfMessages, expectedMessageData)
}
