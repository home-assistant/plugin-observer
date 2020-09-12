package main

import (
	"log"
	"net/http"
	"time"
)

func supervisorPing() bool {
	httpClient := http.Client{
		Timeout: 3 * time.Second,
	}
	response, err := httpClient.Get("http://supervisor/supervisor/ping")
	if err != nil {
		log.Printf("Supervisor ping failed with error %s", err)
		return false
	}

	// Check response
	if response.StatusCode < 300 {
		return true
	}

	log.Printf("Supervisor ping failed with %d", response.StatusCode)
	return false
}
