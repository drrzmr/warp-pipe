package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pagarme/warp-pipe/cmd/reader"
	"github.com/pagarme/warp-pipe/cmd/version"
	"github.com/pagarme/warp-pipe/config"
)

var root *cobra.Command

func init() {

	root = &cobra.Command{
		Use:   config.AppName,
		Short: config.AppShortDescription,
		Long:  "",
	}

	root.AddCommand(version.New())
	root.AddCommand(reader.New())
}

// Execute executes root command
func Execute() error {
	return root.Execute()
}
