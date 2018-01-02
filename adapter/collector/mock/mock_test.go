package mock_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pagarme/warp-pipe/adapter/collector/mock"
	"github.com/pagarme/warp-pipe/pipeline"
)

func TestMockCollector(t *testing.T) {

	collect := func(messageID uint64, publishCh chan<- pipeline.Message) (end bool) {

		publishCh <- pipeline.NewMessage(pipeline.Payload{
			"data": messageID,
		})
		fmt.Println("send message, messageID:", messageID)
		return false
	}

	updateOffset := func(offset uint64) {

		fmt.Println("update offset for message, offset", offset)
	}

	var (
		collector = mock.New(10, collect, updateOffset)
		publishCh = make(chan pipeline.Message)
		offsetCh  = make(chan uint64)
		ctx       = context.Background()
	)

	collector.Init(ctx)
	go collector.Collect(publishCh)
	go collector.UpdateOffset(offsetCh)

	func(publishCh <-chan pipeline.Message, offsetCh chan<- uint64) {
		defer close(offsetCh)

		for msg := range publishCh {
			messageID := msg.Get("data").(uint64)
			fmt.Println("consume message, messageID:", messageID)
			offsetCh <- messageID
		}
	}(publishCh, offsetCh)
}
