package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pagarme/warp-pipe/cmd/dump"
	"github.com/pagarme/warp-pipe/cmd/producer"
	"github.com/pagarme/warp-pipe/cmd/reader"
	"github.com/pagarme/warp-pipe/cmd/version"
	"github.com/pagarme/warp-pipe/config"
	"github.com/pagarme/warp-pipe/lib/log"
)

// Root object
type Root struct {
	conf    *config.Config
	command *cobra.Command
}

// New create a new Root object
func New() *Root {

	var (
		conf    = config.New()
		command *cobra.Command
		stdout  string
		stderr  string
	)

	command = &cobra.Command{
		Use:   config.AppName,
		Short: config.AppShortDescription,
		Long:  "",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			config.String(&conf.Log.Stdout, stdout, config.Default.Log.Stdout)
			config.String(&conf.Log.Stderr, stderr, config.Default.Log.Stderr)

			log.Setup(conf.Log)
		},
	}

	command.AddCommand(version.New())
	command.AddCommand(reader.New(&conf.Reader))
	command.AddCommand(producer.New())
	command.AddCommand(dump.New(conf))

	flags := command.PersistentFlags()
	flags.StringVar(&stdout, "log", config.Default.Log.Stdout, "log stdout")
	flags.StringVar(&stderr, "error", config.Default.Log.Stderr, "log stderr")

	return &Root{
		conf:    conf,
		command: command,
	}
}

// Execute executes root command
func (r *Root) Execute() error {
	return r.command.Execute()
}
