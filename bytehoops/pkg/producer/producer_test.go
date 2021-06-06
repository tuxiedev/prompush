package producer

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	testProducerConfig := Config{
		BootstrapBrokers: []string{"localhost:9092"},
		Topic:            "prometheus",
		ListenAddress:    ":8080",
	}
	router := setupApp(testProducerConfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/write", bytes.NewBufferString("hello world"))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// TODO Add test with headers and consume data from kafka to verify receipt
// TODO Add test without headers and consume data from kafka to verify receipt
