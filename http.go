package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
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

func apiLogs(w http.ResponseWriter, r *http.Request) {
	if !checkNetwork(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Access to logs from %s", r.RemoteAddr)

	err := supervisorLogs(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the content
	w.Header().Add("Content-Type", "text/plain")
}

func apiRestart(w http.ResponseWriter, r *http.Request) {
	if !checkNetwork(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Access to restart from %s", r.RemoteAddr)

	err := supervisorRestart()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the content
	w.WriteHeader(http.StatusOK)
}

type statusData struct {
	On   bool
	Logs string
}

func statusIndex(w http.ResponseWriter, r *http.Request) {
	data := statusData{
		On:   supervisorPing(),
		Logs: "",
	}

	// Set logs
	if data.On {
		var buf bytes.Buffer
		var re = regexp.MustCompile(`\[\d+m`)
		logWriter := bufio.NewWriter(&buf)

		logs := supervisorLogs(logWriter)

		data.Logs = re.ReplaceAllLiteralString(buf.String(), "")

		// Fallback to error
		if data.Logs == "" {
			data.Logs = logs.Error()
		}
	}

	// Render Website
	indexTemplate.Execute(w, data)
}
