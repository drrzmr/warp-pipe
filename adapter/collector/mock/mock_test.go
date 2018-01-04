package mock_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pagarme/warp-pipe/adapter/collector/mock"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

func TestMockCollector(t *testing.T) {

	collect := func(messageID uint64, publishCh chan<- message.Message) (end bool) {

		publishCh <- message.New(message.Payload{
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
		publishCh = make(chan message.Message)
		offsetCh  = make(chan uint64)
		ctx       = context.Background()
	)

	collector.Init(ctx)
	go collector.Collect(publishCh)
	go collector.UpdateOffset(offsetCh)

	func(publishCh <-chan message.Message, offsetCh chan<- uint64) {
		defer close(offsetCh)

		for msg := range publishCh {
			messageID := msg.Get("data").(uint64)
			fmt.Println("consume message, messageID:", messageID)
			offsetCh <- messageID
		}
	}(publishCh, offsetCh)
}
