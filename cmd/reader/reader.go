package reader

import (
	"github.com/pagarme/warp-pipe/config"
	"github.com/pagarme/warp-pipe/reader"
	"github.com/spf13/cobra"
)

// New returns a reader command
func New(configReader *config.Reader) *cobra.Command {

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
