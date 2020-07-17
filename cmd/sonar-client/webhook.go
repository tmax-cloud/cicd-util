package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/cqbqdd11519/cicd-util/pkg/utils"
)

var resultPath string
var outputPath string

func launchWebhook() {
	var err error

	port := WebhookPort
	portStr := os.Getenv("WEBHOOK_PORT")
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			utils.ExitError(log, err, "cannot parse WEBHOOK_PORT")
		}
	}

	resultPath = os.Getenv("SONAR_RESULT_FILE")
	if resultPath == "" {
		resultPath = SonarResultPath
	}

	// Server
	addr := fmt.Sprintf(":%d", port)
	server := http.Server{
		Addr: addr,
	}

	// Router
	log.Info(fmt.Sprintf("Handler set to /"))
	http.HandleFunc("/", handler)

	// Start server
	log.Info(fmt.Sprintf("Server is running on %s", addr))
	if err := server.ListenAndServe(); err != nil {
		log.Error(err, "cannot listen")
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "cannot read body"
		log.Error(err, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Info(fmt.Sprintf("Webhook arrived... Body : %s", string(body)))

	if err := ioutil.WriteFile(resultPath, body, 0777); err != nil {
		msg := "cannot write to file"
		log.Error(err, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintln(w, "")
	if err != nil {
		msg := "cannot respond"
		log.Error(err, "msg")
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	os.Exit(0)
}
