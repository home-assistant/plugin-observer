package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"text/template"

	"github.com/docker/docker/client"
)

var cli *client.Client
var apiKey string
var hassioNetwork *net.IPNet
var indexTemplate *template.Template

func main() {
	var err error

	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	apiKey = os.Getenv("SUPERVISOR_TOKEN")
	_, hassioNetwork, _ = net.ParseCIDR("172.30.32.0/23")
	indexTemplate = template.Must(template.ParseFiles("/usr/share/www/index.html"))

	http.HandleFunc("/", statusIndex)
	http.HandleFunc("/logs", supervisorLogs)
	http.HandleFunc("/restart", supervisorRestart)

	log.Print("Start internal API")
	http.ListenAndServe(":80", nil)
}
