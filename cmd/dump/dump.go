package dump

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jackc/pgx"
	"github.com/pagarme/warp-pipe/lib/namedpipe"
	"github.com/pagarme/warp-pipe/lib/parser/testdecoding"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	postgresCollector "github.com/pagarme/warp-pipe/adapter/collector/postgres"
	testdecodingProcessor "github.com/pagarme/warp-pipe/adapter/processor/testdecoding"
	"github.com/pagarme/warp-pipe/config"
	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/pipeline"
)

var logger *zap.Logger

func init() { log.Register(&logger, "cmd.dump") }

// New returns a dump command
func New(conf *config.Config) (command *cobra.Command) {

	var (
		port     uint16
		host     string
		database string
		slot     string
		output   string
		input    string
	)

	command = &cobra.Command{
		Use:   "dump",
		Short: "Start dump mode",
		Long:  "",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			config.Uint16(&conf.Postgres.Port, port, config.Default.Postgres.Port)
			config.String(&conf.Postgres.Host, host, config.Default.Postgres.Host)
			config.String(&conf.Postgres.Database, database, config.Default.Postgres.Database)
			config.String(&conf.Postgres.Replicate.Slot, slot, config.Default.Postgres.Replicate.Slot)
			config.String(&conf.Cmd.Dump.Stdout, output, config.Default.Cmd.Dump.Stdout)
			config.String(&conf.Cmd.Dump.InputNamedPipe, input, config.Default.Cmd.Dump.Stdout)

			var err error

			// set dump output
			output := os.Stdout
			if conf.Cmd.Dump.Stdout != "stdout" {
				output, err = os.Create(conf.Cmd.Dump.Stdout)
				if err != nil {
					logger.Error("could not create output file", zap.String("name", conf.Cmd.Dump.Stdout))
					return
				}
				defer output.Close()
			}

			// log dump config
			buf, err := json.MarshalIndent(conf, ">", "\t")
			no(err)
			fmt.Fprintf(output, ">>>>>> config begin <<<<<<\n")
			fmt.Fprintf(output, ">%s\n", string(buf))
			fmt.Fprintf(output, ">>>>>> config end <<<<<<\n")

			// set dump input
			pipe := namedpipe.New(conf.Cmd.Dump.InputNamedPipe)
			if err = pipe.Create(); err != nil {
				logger.Error("could not create input named pipe",
					zap.String("name", conf.Cmd.Dump.InputNamedPipe),
				)
				return
			}

			if err = dump(conf, output, pipe); err != nil {
				logger.Error("could not execute dump", zap.Error(err))
			}
		},
	}

	flags := command.Flags()
	flags.Uint16Var(&port, "port", config.Default.Postgres.Port, "postgres port")
	flags.StringVar(&host, "host", config.Default.Postgres.Host, "postgres host")
	flags.StringVar(&database, "database", config.Default.Postgres.Database, "postgres database")
	flags.StringVar(&slot, "slot", config.Default.Postgres.Replicate.Slot, "postgres replicate slot")
	flags.StringVar(&output, "output", config.Default.Cmd.Dump.Stdout, "dump output")
	flags.StringVar(&input, "input", config.Default.Cmd.Dump.InputNamedPipe, "dump input")

	return command
}

func dump(conf *config.Config, output io.Writer, pipe *namedpipe.NamedPipe) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipelineOut, offsetIn, err := createDumpPipeline(ctx, conf.Postgres)
	if err != nil {
		return errors.WithStack(err)
	}

	go handleInput(ctx, offsetIn, pipe)

	for msg := range pipelineOut {
		transaction := msg.Get("transaction").(testdecoding.Transaction)

		var buf []byte

		buf, err = json.MarshalIndent(transaction, "", "\t")
		if err != nil {
			logger.Error("json marshal", zap.Error(err))
			continue
		}

		fmt.Fprintf(output, "%s\n", buf)
	}

	return nil
}

func createDumpPipeline(ctx context.Context, config postgres.Config) (
	out <-chan pipeline.Message, offset chan<- uint64, err error) {

	var (
		run       = pipeline.NewRunner("dump")
		collector = postgresCollector.New(config)
		processor = testdecodingProcessor.New()
	)

	collectorOut, offsetIn, err := run.Collector(ctx, collector)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	processorOut, err := run.Processor(ctx, processor, collectorOut)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return processorOut, offsetIn, nil
}

func handleInput(ctx context.Context, offsetCh chan<- uint64, pipe *namedpipe.NamedPipe) {
	var (
		lineCh, errCh = pipe.Loop(ctx)

		ok     = true
		line   string
		err    error
		offset uint64
	)

	for ok {
		select {
		case line, ok = <-lineCh:
			if !ok {
				break
			}

			offset, err = pgx.ParseLSN(line)
			if err != nil {
				logger.Error("input error", zap.String("line", line), zap.Error(err))
				break
			}
			logger.Info("update consumed", zap.String("lsn", line))
			offsetCh <- offset

		case err, ok = <-errCh:
			if !ok {
				break
			}
			logger.Error("read pipe", zap.Error(err))
		}
	}
}

func no(err error) {
	if err != nil {
		panic(err)
	}
}
