package main

import (
	"context"
	"github.com/cqbqdd11519/cicd-util/pkg/utils"
	regv1 "github.com/tmax-cloud/registry-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func signImage() {
	NAME := "cicd-test"
	ns, err := utils.Namespace()
	if err != nil {
		utils.ExitError(log, err, "cannot get current namespace")
	}

	var img, signer, secretName string
	utils.GetEnvOrDie("IMAGE_URL", &img, log)
	utils.GetEnvOrDie("SIGNER", &signer, log)
	utils.GetEnvOrDie("DOCKER_SECRET_NAME", &secretName, log)

	scheme := runtime.NewScheme()
	if err := regv1.AddToScheme(scheme); err != nil {
		utils.ExitError(log, err, "")
	}

	c, err := utils.Client(client.Options{Scheme: scheme})
	if err != nil {
		utils.ExitError(log, err, "cannot get k8s client")
	}

	req := &regv1.ImageSignRequest{
		ObjectMeta: metav1.ObjectMeta{Name: NAME, Namespace: ns},
		Spec: regv1.ImageSignRequestSpec{
			Image:          img,
			Signer:         signer,
			RegistrySecret: regv1.RegistrySecret{DcjSecretName: secretName},
		},
	}

	if err := c.Create(context.Background(), req); err != nil {
		utils.ExitError(log, err, "")
	}

	// Let's poll...
	for {
		ret := -1
		if err := c.Get(context.Background(), types.NamespacedName{Name: NAME, Namespace: ns}, req); err != nil {
			log.Error(err, "")
		} else if req.Status.ImageSignResponse == nil {
			// Do nothing
		} else if req.Status.ImageSignResponse.Result == regv1.ResponseResultSuccess {
			log.Info("Successfully signed image")
			ret = 0
		} else if req.Status.ImageSignResponse.Result == regv1.ResponseResultFail {
			log.Info("Error while signing image")
			ret = 1
		}

		if ret >= 0 {
			os.Exit(exitWithDelete(c, req, ret))
		}

		time.Sleep(5 * time.Second)
	}
}
