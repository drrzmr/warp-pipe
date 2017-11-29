package mock_test

import (
	"fmt"
	"testing"

	"github.com/pagarme/warp-pipe/adapter/collector/mock"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

func TestMockCollector(t *testing.T) {

	collect := func(c *mock.Collector, messageID uint64, publishCh chan<- message.Message) (end bool) {

		publishCh <- message.New(message.Payload{
			"data": messageID,
		})
		fmt.Println("[mock] send message:", messageID)
		return false
	}

	updateOffset := func(c *mock.Collector, offset uint64) {
		fmt.Println("[mock] update offset for message:", offset)
	}

	mockCollector := mock.New(10, collect, updateOffset)

	publishCh := make(chan message.Message)
	offsetCh := make(chan uint64)

	mockCollector.Init(publishCh, offsetCh)
	go mockCollector.Collect()
	go mockCollector.UpdateOffset()

	func(publishCh <-chan message.Message, offsetCh chan<- uint64) {
		defer close(offsetCh)

		for msg := range publishCh {
			messageID := msg.Get("data").(uint64)
			fmt.Println("[test] consume message:", messageID)
			offsetCh <- messageID
		}
	}(publishCh, offsetCh)
}
