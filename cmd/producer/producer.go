package producer

import (
	"github.com/pagarme/warp-pipe/producer"
	"github.com/spf13/cobra"
)

// New returns a producer command
func New() *cobra.Command {

	return &cobra.Command{
		Use:   "producer",
		Short: "Send messages to Kafka",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

			if err := producer.Run(); err != nil {
				panic(err)
			}
		},
	}
}
