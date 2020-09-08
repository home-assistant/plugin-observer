package main

import (
	"net/http"

	"github.com/docker/docker/client"
)

var cli *client.Client

func main() {
	var err error

	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/logs", supervisorLogs)
	http.ListenAndServe(":80", nil)
}
