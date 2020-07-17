package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/tidwall/gjson"

	"github.com/cqbqdd11519/cicd-util/pkg/sonar"
	"github.com/cqbqdd11519/cicd-util/pkg/utils"
)

func post() {
	sonarUrl := os.Getenv("SONAR_URL")
	if sonarUrl == "" {
		utils.ExitError(log, fmt.Errorf("SONAR_URL not set"), "environment doesn't meet condition")
	}

	sonarToken := os.Getenv("SONAR_TOKEN")
	if sonarToken == "" {
		utils.ExitError(log, fmt.Errorf("SONAR_TOKEN not set"), "environment doesn't meet condition")
	}

	sonarResultPath := os.Getenv("SONAR_RESULT_FILE")
	if sonarResultPath == "" {
		sonarResultPath = SonarResultPath
	}

	outputPath := os.Getenv("SONAR_RESULT_DEST")
	if outputPath == "" {
		outputPath = SonarOutputPath
	}

	webhookKeyPath := os.Getenv("SONAR_WEBHOOK_KEY_FILE")
	if webhookKeyPath == "" {
		webhookKeyPath = WebhookKeyPath
	}
	webhookKey, err := ioutil.ReadFile(webhookKeyPath)
	if err != nil {
		utils.ExitError(log, err, "cannot read webhook key file")
	}

	// Wait until webhook result arrives
	if !utils.FileExists(sonarResultPath) {
		w := watcher.New()
		w.SetMaxEvents(1)
		w.FilterOps(watcher.Create)

		dir, err := filepath.Abs(filepath.Dir(sonarResultPath))
		if err != nil {
			utils.ExitError(log, err, "cannot get dirname of sonar result file")
		}
		fileName := filepath.Base(sonarResultPath)

		if err := w.Add(dir); err != nil {
			utils.ExitError(log, err, "cannot watch directory")
		}

		done := make(chan error)
		go func() {
			for {
				select {
				case event := <-w.Event:
					if event.Name() == fileName {
						done <- nil
					}
				case err := <-w.Error:
					done <- err
				case <-w.Closed:
					done <- nil
				}
			}
		}()

		go func() {
			if err := w.Start(time.Millisecond * 100); err != nil {
				utils.ExitError(log, err, "cannot start watcher")
			}
		}()

		if err := <-done; err != nil {
			utils.ExitError(log, err, "error occurred at watching file")
		}
	}

	// Delete Webhook
	if err := sonar.DeleteWebhook(sonarUrl, sonarToken, string(webhookKey)); err != nil {
		utils.ExitError(log, err, "cannot delete webhook")
	}

	// Decide OK
	result, err := ioutil.ReadFile(sonarResultPath)
	if err != nil {
		utils.ExitError(log, err, "cannot read result file")
	}

	if err := ioutil.WriteFile(outputPath, result, 0777); err != nil {
		utils.ExitError(log, err, "cannot write result")
	}

	qualityStatus := gjson.Get(string(result), "qualityGate.status")
	if !qualityStatus.Exists() {
		utils.ExitError(log, err, "there is no qualityGate.status value")
	}

	log.Info(fmt.Sprintf("Status: %s", qualityStatus.String()))
	if qualityStatus.String() == "OK" {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
