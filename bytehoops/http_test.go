package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupApp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/write", bytes.NewBufferString("hello world"))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// TODO Add test with headers and consume data from kafka to verify receipt
// TODO Add test without headers and consume data from kafka to verify receipt
