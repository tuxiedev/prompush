package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuxiedev/prompush/bytehoops/pkg/consumer"
)

var (
	consumerGroupName string
	sinkEndpoint      string
)

func init() {
	rootCmd.AddCommand(consumerCmd)
	consumerCmd.Flags().StringVarP(&consumerGroupName, "group-name", "g",
		"bytehoops", "Name of the Kafka consumer group")
	consumerCmd.Flags().StringVarP(&sinkEndpoint, "sink-endpoint", "s",
		"http://localhost:9009/api/v1/push", "Path to remote http endpoint")
}

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consumes byte arrays from Kafka and posts messages to an API endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		consumer.RunConsumer(consumer.Config{
			BootstrapBrokers:  bootstrapBrokers,
			Topic:             topic,
			ConsumerGroupName: consumerGroupName,
			SinkEndpoint:      sinkEndpoint,
		})
	},
}
