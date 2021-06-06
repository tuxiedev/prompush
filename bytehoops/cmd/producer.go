package cmd

import (
	"github.com/tuxiedev/prompush/bytehoops/pkg/producer"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(producerCmd)
	producerCmd.Flags().StringP("listen-address", "l", ":8080", "Listening address for the API server to receive events")
}

var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "Produces a byte array published to this http endpoint to Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		listenAddress, err := cmd.Flags().GetString("listen-address")
		if err != nil {
			log.Println("Could not get listen address, falling back to :8080")
			listenAddress = ":8080"
		}
		producer.RunProducer(producer.Config{
			BootstrapBrokers: bootstrapBrokers,
			Topic:            topic,
			ListenAddress:    listenAddress,
		})
	},
}
