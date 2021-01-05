package utils

import (
	"fmt"
	"github.com/go-logr/logr"
	"os"
)

func GetEnv(key string, value *string) error {
	*value = os.Getenv(key)
	if *value == "" {
		return fmt.Errorf("env %s is not found", key)
	}
	return nil
}

func GetEnvOrDie(key string, value *string, log logr.Logger) {
	if err := GetEnv(key, value); err != nil {
		ExitError(log, err, "")
	}
}
