package consumer

type Config struct {
	BootstrapBrokers  []string
	Topic             string
	ConsumerGroupName string
	SinkEndpoint      string
}
