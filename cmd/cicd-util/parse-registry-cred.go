package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tmax-cloud/cicd-util/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func parseRegistryCred() {
	// Check env. var.s
	targetFilePath := os.Getenv("TARGET_FILE")
	if targetFilePath == "" {
		utils.ExitError(log, fmt.Errorf("environment not given"), "TARGET_FILE should be given")
	}
	secretName := os.Getenv("SECRET_NAME")
	if secretName == "" {
		log.Info("no SECRET_NAME is given... skipping making credential")
		writeCredToFile(targetFilePath, "")
		return
	}
	imageUrlFilePath := os.Getenv("IMAGE_URL_FILE")
	if imageUrlFilePath == "" {
		utils.ExitError(log, fmt.Errorf("environment not given"), "IMAGE_URL_FILE should be given")
	}

	// Read image url from file
	imageUrl, err := ioutil.ReadFile(imageUrlFilePath)
	if err != nil {
		utils.ExitError(log, err, "cannot read IMAGE_URL_FILE")
	}

	// Get registry base url
	registry := strings.Split(string(imageUrl), "/")[0]

	ns, err := utils.Namespace()
	if err != nil {
		utils.ExitError(log, err, "cannot get current namespace")
	}

	c, err := utils.Client(client.Options{})
	if err != nil {
		utils.ExitError(log, err, "cannot get client")
	}

	secret := &corev1.Secret{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: ns}, secret); err != nil {
		utils.ExitError(log, err, "cannot get secret")
	}

	credBytes, exist := secret.Data[corev1.DockerConfigJsonKey]
	if !exist {
		utils.ExitError(log, fmt.Errorf("no exptected secret key"), ".dockerconfigjson should exist in secret")
	}
	credStr := string(credBytes)

	key := fmt.Sprintf("auths.%s", escapeDot(registry))
	cred := gjson.Get(credStr, key)
	if !cred.Exists() {
		log.Info(fmt.Sprintf("There is no key %s in JSON", key))
		writeCredToFile(targetFilePath, "")
		os.Exit(0)
	}

	auth := cred.Get("auth")
	if !auth.Exists() {
		utils.ExitError(log, fmt.Errorf("key auth does not exist in found cred json object"), "no auth found in secret")
	}

	writeCredToFile(targetFilePath, auth.String())
}

func escapeDot(src string) string {
	return strings.ReplaceAll(src, ".", "\\.")
}

func writeCredToFile(filePath, cred string) {
	log.Info(fmt.Sprintf("Credential : %s", cred))
	if err := ioutil.WriteFile(filePath, []byte(cred), 0777); err != nil {
		utils.ExitError(log, err, "cannot write to file")
	}
}
