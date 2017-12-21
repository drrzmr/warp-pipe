package mock_test

import (
	"testing"

	"go.uber.org/zap"

	"github.com/pagarme/warp-pipe/adapter/collector/mock"
	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/pipeline/message"
)

var logger = log.Development("test")

func TestMockCollector(t *testing.T) {

	collect := func(c *mock.Collector, messageID uint64, publishCh chan<- message.Message) (end bool) {

		publishCh <- message.New(message.Payload{
			"data": messageID,
		})
		logger.Debug("send message", zap.Uint64("messageID", messageID))
		return false
	}

	updateOffset := func(c *mock.Collector, offset uint64) {

		logger.Debug("update offset for message", zap.Uint64("offset", offset))
	}

	var (
		collector = mock.New(10, collect, updateOffset)
		publishCh = make(chan message.Message)
		offsetCh  = make(chan uint64)
	)

	collector.Init()
	go collector.Collect(publishCh)
	go collector.UpdateOffset(offsetCh)

	func(publishCh <-chan message.Message, offsetCh chan<- uint64) {
		defer close(offsetCh)

		for msg := range publishCh {
			messageID := msg.Get("data").(uint64)
			logger.Debug("consume message", zap.Uint64("messageID", messageID))
			offsetCh <- messageID
		}
	}(publishCh, offsetCh)
}
