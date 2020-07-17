package utils

import (
	"github.com/go-logr/logr"
	"os"
)

func ExitError(log logr.Logger, err error, msg string) {
	log.Error(err, msg)
	os.Exit(1)
}
