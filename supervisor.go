package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

type SupervisorResponse struct {
	Result  string                 `json:"result"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type SupervisorInfo struct {
	Healthy   bool `json:"healthy"`
	Supported bool `json:"supported"`
}

type ResolutionInfo struct {
	Unhealthy   []string `json:"unhealthy"`
	Unsupported []string `json:"unsupported"`
}

type SupervisorPing struct {
	Connected bool
	State     string
}

func supervisorApiProxy(path string) (SupervisorResponse, error) {
	var jsonResponse SupervisorResponse
	request, _ := http.NewRequest("GET", fmt.Sprintf("http://supervisor/%s", path), nil)
	request.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", os.Getenv("SUPERVISOR_TOKEN"))},
	}

	response, err := httpClient.Do(request)
	if err != nil {
		log.Printf("Supervisor API call failed with error %s", err)
		return jsonResponse, err
	}

	if response.StatusCode >= 300 && response.StatusCode != 400 {
		log.Printf("Supervisor API call failed with status code %v", response.StatusCode)
		return jsonResponse, fmt.Errorf("call to Supervisor API failed with status code %v", response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return jsonResponse, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	err = json.Unmarshal([]byte(bodyBytes), &jsonResponse)
	if err != nil {
		log.Printf("json unmarshalling failed with error %s", err)
		return jsonResponse, err
	}

	if response.StatusCode == 400 {
		return jsonResponse, errors.New("call to Supervisor API failed with status code 400")
	}

	return jsonResponse, err
}

func supervisorPing() SupervisorPing {
	supervisorPingData := SupervisorPing{
		Connected: true,
	}
	data, err := supervisorApiProxy("supervisor/ping")
	if err != nil {
		log.Printf("Supervisor ping failed with error %s", err)
		supervisorPingData.Connected = false
		if strings.HasPrefix(data.Message, "System is not ready with state:") {
			// This is an API error, but we got a proper response so we accept it
			supervisorPingData.Connected = true
			supervisorPingData.State = strings.ReplaceAll(data.Message, "System is not ready with state: ", "")
		}
	}
	return supervisorPingData
}

func getSupervisorInfo() (SupervisorInfo, error) {
	var supervisorInfo SupervisorInfo
	response, err := supervisorApiProxy("supervisor/info")
	if err != nil {
		log.Printf("Supervisor API call failed with error %s", err)
		return supervisorInfo, err
	}

	jsonData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(jsonData, &supervisorInfo)
	if err != nil {
		log.Printf("json unmarshalling failed with error %s", err)
		return supervisorInfo, err
	}

	return supervisorInfo, nil
}

func getResolutionInfo() (ResolutionInfo, error) {
	var resolutionInfo ResolutionInfo
	response, err := supervisorApiProxy("resolution/info")
	if err != nil {
		log.Printf("Supervisor API call failed with error %s", err)
		return resolutionInfo, err
	}

	jsonData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(jsonData, &resolutionInfo)
	if err != nil {
		log.Printf("json unmarshalling failed with error %s", err)
		return resolutionInfo, err
	}

	return resolutionInfo, nil
}

func supervisorLogs(w io.Writer) error {
	// Read logs from container
	reader, err := cli.ContainerLogs(context.Background(), "hassio_supervisor", container.LogsOptions{
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
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("error closing reader: %v", err)
		}
	}()



	// Return the content
	_, err = stdcopy.StdCopy(w, w, reader)
	return err
}
