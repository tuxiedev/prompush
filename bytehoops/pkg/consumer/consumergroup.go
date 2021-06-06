package consumer

import (
	"bytes"
	"github.com/Shopify/sarama"
	"github.com/gojek/heimdall/v7/httpclient"
	"log"
	"strings"
)

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

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		headers := getHttpHeaderFromKafkaHeaders(message.Headers)
		_, err := consumer.httpClient.Post(consumer.ingestUrl, bytes.NewBuffer(message.Value), headers)
		if err != nil {
			log.Println("Error with request", err)
		}
		//TODO deal with errors:
		// Note: 429 error occurs when ingest rate limit exceeds
		session.MarkMessage(message, "")
	}

	return nil
}
