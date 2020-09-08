package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
)

func supervisorLogs(w http.ResponseWriter, r *http.Request) {
	reader, _ := cli.ContainerLogs(context.Background(), "hassio_supervisor", types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	defer reader.Close()
	content, _ := ioutil.ReadAll(reader)

	w.Write(content)
}
