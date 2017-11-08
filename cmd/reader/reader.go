package reader

import (
	"os"

	"github.com/pagarme/warp-pipe/reader"
	"github.com/spf13/cobra"
)

// New returns a reader command
func New() *cobra.Command {

	return &cobra.Command{
		Use:   "reader",
		Short: "Start reader mode",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

			reader.Run(os.Stdin, os.Stdout)
		},
	}
}
