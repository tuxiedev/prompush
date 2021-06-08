package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/plugins"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func createHttpClient() *httpclient.Client {
	//TODO Following options should be configurable by user
	initialInterval := 2 * time.Second
	maxTimeout := 120 * time.Second
	maximumJitterInterval := 5 * time.Second
	backoff := heimdall.NewExponentialBackoff(initialInterval, maxTimeout, 2, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	timeout := 1000 * time.Millisecond
	// Create a new client, sets the retry mechanism, and the number of times you would like to retry
	return httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(4),
	)
}

func RunConsumer(config Config) {
	client := createHttpClient()
	requestLogger := plugins.NewRequestLogger(nil, nil)
	client.AddPlugin(requestLogger)
	consumer := Consumer{
		ready:      make(chan bool),
		httpClient: *client,
		ingestUrl:  config.SinkEndpoint,
	}
	saramaConfig := sarama.NewConfig()
	ctx, cancel := context.WithCancel(context.Background())
	consumerClient, err := sarama.NewConsumerGroup(config.BootstrapBrokers, config.ConsumerGroupName, saramaConfig)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := consumerClient.Consume(ctx, []string{config.Topic}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = consumerClient.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
