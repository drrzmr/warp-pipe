package testdecoding

import (
	"context"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/parser/testdecoding"
	"github.com/pagarme/warp-pipe/pipeline"
	"go.uber.org/zap"
)

// Processor object
type Processor struct {
	pipeline.Processor
	ctx context.Context
	in  <-chan pipeline.Message
}

var logger *zap.Logger

func init() { log.Register(&logger, "adapter.processor.testdecoding") }

// New create a new processor
func New() *Processor {

	return &Processor{
		ctx: nil,
		in:  nil,
	}
}

// Init implement interface
func (p *Processor) Init(ctx context.Context, in <-chan pipeline.Message) (err error) {
	logger.Debug("--> Init()")
	defer logger.Debug("<-- Init()")

	if p.isInitialized() {
		logger.Warn("already initialized")
		return nil
	}

	p.ctx = ctx
	p.in = in

	return nil
}

// Process implement interface
func (p *Processor) Process(out chan<- pipeline.Message) {
	defer close(out)

	parser := testdecoding.NewParser(func(transaction testdecoding.Transaction) {
		out <- pipeline.NewMessage(pipeline.Payload{
			"transaction": transaction,
		})
	})

	for msg := range p.in {

		walData := string(msg.Get("WalData").([]byte))

		logger.Debug("event", zap.String("message", walData))

		err := parser.Parse(walData)
		if err != nil {
			panic(err)
		}

	}
}

func (p *Processor) isInitialized() (initialized bool) { return p.ctx != nil && p.in != nil }
