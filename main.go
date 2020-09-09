package main

import (
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/client"
)

var cli *client.Client
var apiKey string

func main() {
	var err error

	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	apiKey = os.Getenv("OBSERVER_TOKEN")

	http.HandleFunc("/logs", supervisorLogs)
	http.HandleFunc("/restart", supervisorRestart)

	log.Print("Start internal API")
	http.ListenAndServe(":80", nil)
}
