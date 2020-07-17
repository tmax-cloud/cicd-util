package main

import (
	"fmt"
	"os"

	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cqbqdd11519/cicd-util/pkg/utils"
)

var log = logf.Log.WithName("cicd-util")

func main() {
	logf.SetLogger(zap.Logger())

	args := os.Args[1:]

	if len(args) != 1 {
		utils.ExitError(log, fmt.Errorf("there should be only one argument"), "doesn't meet argument condition")
	}

	switch args[0] {
	case "pre":
		pre()
	case "post":
		post()
	case "webhook":
		launchWebhook()
	default:
		utils.ExitError(log, fmt.Errorf("command should be one if [pre|post|webhook]"), "not supported")
	}
}
