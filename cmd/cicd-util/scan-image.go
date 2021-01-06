package main

import (
	"context"
	"fmt"
	"github.com/cqbqdd11519/cicd-util/pkg/utils"
	scanv1 "github.com/tmax-cloud/image-scanning-operator/api/v1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"time"
)

func scanImage() {
	NAME := "cicd-test"
	ns, err := utils.Namespace()
	if err != nil {
		utils.ExitError(log, err, "cannot get current namespace")
	}

	var img, thresholdStr string
	utils.GetEnvOrDie("IMAGE_URL", &img, log)
	utils.GetEnvOrDie("THRESHOLD", &thresholdStr, log)

	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		utils.ExitError(log, err, "")
	}

	scheme := runtime.NewScheme()
	if err := scanv1.AddToScheme(scheme); err != nil {
		utils.ExitError(log, err, "")
	}

	c, err := utils.Client(client.Options{Scheme: scheme})
	if err != nil {
		utils.ExitError(log, err, "cannot get k8s client")
	}

	req := &scanv1.ImageScanning{
		ObjectMeta: metav1.ObjectMeta{Name: NAME, Namespace: ns},
		Spec: scanv1.ImageScanningSpec{
			ImageUrl:    img,
			ForceNonSSL: true,
			Insecure:    true,
			Webhook:     true,
		},
	}

	if err := c.Create(context.Background(), req); err != nil {
		utils.ExitError(log, err, "")
	}

	// Let's poll...
	for {
		ret := -1
		if err := c.Get(context.Background(), types.NamespacedName{Name: NAME, Namespace: ns}, req); err != nil {
			fmt.Println(err.Error())
		} else {
			printScanStatus(req)
			switch req.Status.Status {
			case scanv1.ScanningSuccess:
				total := 0
				for k, v := range req.Status.Summary {
					if strings.ToLower(k) == "negligible" {
						continue
					}
					total += v
				}

				if total >= threshold {
					fmt.Printf("The number of vulnerabilities (%d) is greater than threshold (%d)\n", total, threshold)
					ret = 1
				} else {
					fmt.Printf("The number of vulnerabilities (%d) is less than threshold (%d)\n", total, threshold)
					ret = 0
				}
			case scanv1.ScanningError:
				fmt.Println("Error while scanning image")
				ret = 1
			}
		}

		if ret >= 0 {
			os.Exit(ret)
		}

		time.Sleep(5 * time.Second)
	}
}

func printScanStatus(req *scanv1.ImageScanning) {
	b, err := yaml.Marshal(req.Status)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("RESULT:")
	fmt.Println(string(b))
}
