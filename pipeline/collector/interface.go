package collector

import "github.com/pagarme/warp-pipe/pipeline/message"

// Collector interface
type Collector interface {
	Init(publishCh chan<- message.Message, updateOffsetCh <-chan uint64) (err error)
	Collect()
	UpdateOffset()
}
