package collector

import (
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/pipeline/message"
)

// Stage object
type Stage struct {
	collector Collector
}

// NewStage return stage
func NewStage(collector Collector) (stage *Stage) {

	return &Stage{
		collector: collector,
	}
}

// Run start stage
func (stage *Stage) Run() (publishCh <-chan message.Message, offsetCh chan<- uint64, err error) {

	publish := make(chan message.Message)
	offset := make(chan uint64)

	if err = stage.collector.Init(publish, offset); err != nil {
		return nil, nil, errors.Wrap(err, "Could not initialize stage collector")
	}

	go stage.collector.Collect()
	go stage.collector.UpdateOffset()

	return publish, offset, nil
}
