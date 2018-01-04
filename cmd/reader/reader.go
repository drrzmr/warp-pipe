package reader

import (
	"github.com/spf13/cobra"

	"github.com/pagarme/warp-pipe/lib/snippet/reader"
)

// New returns a reader command
func New(configReader *reader.Config) *cobra.Command {

	return &cobra.Command{
		Use:   "reader",
		Short: "Start reader mode",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

			in := configReader.InputStream
			out := configReader.OutputStream

			if err := reader.Run(in, out); err != nil {
				panic(err)
			}
		},
	}
}
