package producer

type Config struct {
	BootstrapBrokers []string
	Topic            string
	ListenAddress    string
}
