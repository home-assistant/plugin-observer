package main

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

func supervisorPing() bool {
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

func supervisorLogs(w io.Writer) error {
	// Read logs from container
	reader, err := cli.ContainerLogs(context.Background(), "hassio_supervisor", types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Timestamps: false,
		Tail:       "all",
	})

	if err != nil {
		log.Printf("Can't get supervisor logs %s", err)
		return err
	}
	defer reader.Close()

	// Return the content
	_, err = stdcopy.StdCopy(w, w, reader)
	return err
}

func supervisorRestart() error {
	// Read logs from container
	err := cli.ContainerStop(context.Background(), "hassio_supervisor", nil)
	if err != nil {
		log.Printf("Can't stop supervisor: %s", err)
		return err
	}

	return nil
}
