package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

func checkNetwork(r *http.Request) bool {
	remote := strings.Split(r.RemoteAddr, ":")[0]
	remoteIP := net.ParseIP(remote)

	// Is in network
	if !hassioNetwork.Contains(remoteIP) {
		log.Printf("Access not allow from %s", remote)
		return false
	}

	// If supervisor is down
	if !supervisorPing() {
		log.Printf("API is disabled / Supervisor is running - %s", remote)
		return true
	}

	return false
}

func supervisorLogs(w http.ResponseWriter, r *http.Request) {
	if !checkNetwork(r) {
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
		Follow:     false,
		Timestamps: false,
		Tail:       "all",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	// Return the content
	w.Header().Add("Content-Type", "text/plain")
	stdcopy.StdCopy(w, w, reader)
}

func supervisorRestart(w http.ResponseWriter, r *http.Request) {
	if !checkNetwork(r) {
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

type statusData struct {
	On bool
}

func statusIndex(w http.ResponseWriter, r *http.Request) {
	data := statusData{
		On: supervisorPing(),
	}

	// Render Website
	indexTemplate.Execute(w, data)
}
