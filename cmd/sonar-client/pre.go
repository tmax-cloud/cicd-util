package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cqbqdd11519/cicd-util/pkg/sonar"
	"github.com/cqbqdd11519/cicd-util/pkg/utils"
)

func pre() {
	sonarUrl := os.Getenv("SONAR_URL")
	if sonarUrl == "" {
		utils.ExitError(log, fmt.Errorf("SONAR_URL not set"), "environment doesn't meet condition")
	}

	sonarToken := os.Getenv("SONAR_TOKEN")
	if sonarToken == "" {
		utils.ExitError(log, fmt.Errorf("SONAR_TOKEN not set"), "environment doesn't meet condition")
	}

	sonarProjectId := os.Getenv("SONAR_PROJECT_ID")
	if sonarProjectId == "" {
		utils.ExitError(log, fmt.Errorf("SONAR_PROJECT_ID not set"), "environment doesn't meet condition")
	}

	// Register webhook
	key, err := sonar.RegisterWebhook(sonarUrl, sonarToken, sonarProjectId, WebhookPort)
	if err != nil {
		utils.ExitError(log, err, "cannot register webhook")
	}

	// Save to file
	path := os.Getenv("SONAR_WEBHOOK_KEY_FILE")
	if path == "" {
		path = WebhookKeyPath
	}
	if err := ioutil.WriteFile(path, []byte(key), 0777); err != nil {
		utils.ExitError(log, err, "cannot write webhook key to file")
	}
}
