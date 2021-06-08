package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string
var bootstrapBrokers []string
var topic string

var rootCmd = &cobra.Command{
	Use:   "bytehoops",
	Short: "Bytehoops is a simple way to transport byte array between HTTP endpoints and Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		// noop
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// TODO take all kafka configs from a file instead as there are too many
	rootCmd.PersistentFlags().StringSliceVarP(&bootstrapBrokers, "bootstrap-servers", "b", []string{"localhost:9092"}, "Comma separated list of bootstrap servers")
	rootCmd.PersistentFlags().StringVarP(&topic, "topic", "t", "prometheus", "Topic to consume or send data to")
}
