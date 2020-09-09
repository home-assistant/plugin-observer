package main

import (
	"net/http"
	"os"

	"github.com/docker/docker/client"
)

var cli *client.Client
var apiKey string

func main() {
	var err error

	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	apiKey = os.Getenv("OBSERVER_TOKEN")

	http.HandleFunc("/logs", supervisorLogs)
	http.ListenAndServe(":80", nil)
}
