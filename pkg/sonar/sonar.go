package sonar

import (
	"encoding/base64"
	"fmt"
	"github.com/tidwall/gjson"
	"net/url"
	"path"

	"github.com/cqbqdd11519/cicd-util/pkg/utils"
)

func RegisterWebhook(sonarUrl, sonarToken, projectId string, port int) (string, error) {
	podIp, err := utils.PodIp()
	if err != nil {
		return "", err
	}

	webhookUrl := fmt.Sprintf("http://%s:%d", podIp, port)

	header := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(sonarToken+":"))),
	}
	param := map[string]string{
		"project": projectId,
		"name":    fmt.Sprintf("%s-webhook", projectId),
		"url":     webhookUrl,
	}

	addr, err := url.Parse(sonarUrl)
	if err != nil {
		return "", err
	}
	addr.Path = path.Join(addr.Path, "/api/webhooks/create")
	fmt.Printf("Requesting... %s\n", addr.String())
	result, err := utils.PostReq(addr.String(), header, param)
	if err != nil {
		return "", err
	}

	key := gjson.Get(result, "webhook.key")
	if key.Exists() {
		return key.String(), nil
	} else {
		return "", fmt.Errorf("no webhook.key in json, data: %s", result)
	}
}

func DeleteWebhook(sonarUrl, sonarToken, webhookKey string) error {
	header := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(sonarToken+":"))),
	}
	param := map[string]string{
		"webhook": webhookKey,
	}

	addr, err := url.Parse(sonarUrl)
	if err != nil {
		return err
	}
	addr.Path = path.Join(addr.Path, "/api/webhooks/delete")
	fmt.Printf("Requesting... %s\n", addr.String())

	_, err = utils.PostReq(addr.String(), header, param)
	if err != nil {
		return err
	}

	return nil
}
