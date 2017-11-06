package version

import (
	"github.com/spf13/cobra"

	"github.com/pagarme/warp-pipe/config"
)

// New returns a version command
func New() *cobra.Command {

	return &cobra.Command{
		Use:   "version",
		Short: "Show the version",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("version:", config.AppVersion)
		},
	}
}
