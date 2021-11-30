package hxxp

import (
	"net/http"
	"time"
)

func NewHTTPClient() *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return client
}

