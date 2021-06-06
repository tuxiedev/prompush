package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func getKafkaHeadersFromHttpHeaders(httpHeaders map[string][]string) []sarama.RecordHeader {
	var headers []sarama.RecordHeader
	for key, value := range httpHeaders {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(strings.Join(value, ",")),
		})
	}
	return headers
}

func setupRouter(producer sarama.SyncProducer, topic string) *gin.Engine {
	r := gin.Default()

	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.Any("/v1/write", func(context *gin.Context) {
		body, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			context.AbortWithStatus(400)
		}
		_, _, producerError := producer.SendMessage(&sarama.ProducerMessage{
			Topic:   topic,
			Value:   sarama.ByteEncoder(body),
			Headers: getKafkaHeadersFromHttpHeaders(context.Request.Header),
		})
		if producerError != nil {
			context.AbortWithStatus(500)
		}
		context.String(200, string(body))
	})
	return r
}

func instantiateKafkaProducer() sarama.SyncProducer {
	brokerList := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Cannot instantiate Kafka producer", err)
	}
	return producer
}

func setupApp() *gin.Engine {
	producer := instantiateKafkaProducer()
	topic := "prometheus"
	return setupRouter(producer, topic)
}

type AppConfig struct {
	// TODO Allow ability to configure various components of the app
}

func main() {
	err := setupApp().Run(":8080")
	if err != nil {
		log.Fatalln("Failed to start http server")
	}
}
