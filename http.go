package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
)

func checkAccessKey(r *http.Request) bool {
	token := r.Header.Get("X-Observer-Token")

	// Check api key
	if token != apiKey {
		log.Printf("Unauthorized API access %s", r.RemoteAddr)
		return false
	}
	return true
}

func supervisorLogs(w http.ResponseWriter, r *http.Request) {
	if !checkAccessKey(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Access to logs from %s", r.RemoteAddr)

	// Read logs from container
	reader, err := cli.ContainerLogs(context.Background(), "hassio_supervisor", types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	// Return the content
	content, _ := ioutil.ReadAll(reader)
	w.Write(content)
}

func supervisorRestart(w http.ResponseWriter, r *http.Request) {
	if !checkAccessKey(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Access to restart from %s", r.RemoteAddr)

	// Read logs from container
	err := cli.ContainerStop(context.Background(), "hassio_supervisor", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the content
	w.WriteHeader(http.StatusOK)
}
