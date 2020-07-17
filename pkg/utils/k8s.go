package utils

import (
	"io/ioutil"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	NamespaceFilePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func Namespace() (string, error) {
	ns, err := ioutil.ReadFile(NamespaceFilePath)
	if err != nil {
		envNs := os.Getenv("NAMESPACE")
		if envNs != "" {
			return envNs, nil
		}
		return "", err
	}
	return string(ns), nil
}

func Client(options client.Options) (client.Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	c, err := client.New(cfg, options)
	if err != nil {
		return nil, err
	}
	return c, nil
}
