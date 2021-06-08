package consumer

import (
	"bytes"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/avast/retry-go"
	"github.com/gojek/heimdall/v7/httpclient"
	"log"
	"strings"
	"time"
)

var tooManyRequests = errors.New("Retriable error due to too many requests")

type Consumer struct {
	ready      chan bool
	httpClient httpclient.Client
	ingestUrl  string
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func getHttpHeaderFromKafkaHeaders(kafkaHeaders []*sarama.RecordHeader) map[string][]string {
	var returnMap = map[string][]string{}
	for _, header := range kafkaHeaders {
		returnMap[string(header.Key)] = strings.Split(string(header.Value), ",")
	}
	return returnMap
}

func (consumer *Consumer) processMessage(message *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
	headers := getHttpHeaderFromKafkaHeaders(message.Headers)
	retry.Do(
		func() error {
			resp, err := consumer.httpClient.Post(consumer.ingestUrl, bytes.NewBuffer(message.Value), headers)
			if err != nil {
				log.Println("Error with request", err)
			}
			if resp.StatusCode == 429 || resp.StatusCode == 500 {
				log.Println("Too many requests error occured, retrying")
				return tooManyRequests
			}
			return nil
		},
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			if err == tooManyRequests {
				log.Println("Retrying after backoff on too many requests")
				return retry.BackOffDelay(n, err, config)
			} else {
				return retry.DefaultDelay
			}
		}),
	)
	session.MarkMessage(message, "")
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		consumer.processMessage(message, session)
	}
	return nil
}
