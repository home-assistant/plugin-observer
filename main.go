package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/docker/docker/client"
)

var cli *client.Client
var apiKey string //nolint
var hassioNetwork *net.IPNet
var indexTemplate *template.Template
var wwwRoot string
var development bool
var httpClient http.Client

func main() {
	var err error

	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Generate Network check
	networkMask := "172.30.32.0/23"
	if env := os.Getenv("NETWORK_MASK"); env != "" {
		networkMask = env
	}
	_, hassioNetwork, _ = net.ParseCIDR(networkMask)

	// system settings
	apiKey = os.Getenv("SUPERVISOR_TOKEN")
	development = (os.Getenv("DEVELOPMENT") == "True")
	httpClient = http.Client{Timeout: 3 * time.Second}

	if development {
		wwwRoot = "./rootfs/usr/share/www/"
	} else {
		wwwRoot = "/usr/share/www/"
	}

	indexTemplate = template.Must(template.ParseFiles(wwwRoot + "/index.html"))

	// API setup
	http.HandleFunc("/", statusIndex)
	http.HandleFunc("/ping", apiPing)
	http.HandleFunc("/logs", apiLogs)

	// Serve static help files
	staticFiles := http.FileServer(http.Dir(wwwRoot))
	http.Handle("/static/styles.css", staticFiles)
	http.Handle("/static/scripts.js", staticFiles)
	http.Handle("/static/tsparticles.min.js", staticFiles)
	http.Handle("/static/Roboto-Regular.woff2", staticFiles)
	

	log.Print("Start webserver on http://0.0.0.0:80")
	http.ListenAndServe(":80", nil)
}
