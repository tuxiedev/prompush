package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/hashicorp/go-uuid"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func createHttpClient() *httpclient.Client {
	backoffInterval := 2 * time.Millisecond
	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval := 5 * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

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

func main() {
	client := createHttpClient()
	consumer := Consumer{
		ready:      make(chan bool),
		httpClient: *client,
		ingestUrl:  "http://localhost:9009/api/v1/push",
	}
	config := sarama.NewConfig()
	brokers := []string{"localhost:9092"}
	topics := []string{"prometheus"}
	group, _ := uuid.GenerateUUID()
	ctx, cancel := context.WithCancel(context.Background())
	consumerClient, err := sarama.NewConsumerGroup(brokers, group, config)
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
			if err := consumerClient.Consume(ctx, topics, &consumer); err != nil {
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
