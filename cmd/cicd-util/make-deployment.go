package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	"github.com/ghodss/yaml"

	"github.com/tmax-cloud/cicd-util/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

const (
	TargetFilePath = "/generate/deployment.yaml"
	DepSpecKey     = "deploy-spec.yaml"
	LabelKey       = "app"
)

func makeDeployment() {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		utils.ExitError(log, fmt.Errorf("environment not given"), "APP_NAME should be given")
	}
	imageUrl := os.Getenv("IMAGE_URL")
	if imageUrl == "" {
		utils.ExitError(log, fmt.Errorf("environment not given"), "IMAGE_URL should be given")
	}

	ns, err := utils.Namespace()
	if err != nil {
		utils.ExitError(log, err, "cannot get current namespace")
	}

	depEnvJsonRaw := os.Getenv("DEPLOY_ENV_JSON")
	configMapName := os.Getenv("CONFIGMAP_NAME")

	// K8s client
	c, err := utils.Client(client.Options{})
	if err != nil {
		utils.ExitError(log, err, "cannot get k8s client")
	}

	// Deployment spec that would be merged into dep
	depSpec := &appsv1.Deployment{}
	cm := &corev1.ConfigMap{}
	if configMapName != "" {
		if err := c.Get(context.TODO(), types.NamespacedName{Name: configMapName, Namespace: ns}, cm); err != nil {
			utils.ExitError(log, err, "cannot get configMap")
		}
		depSpecStr, exist := cm.Data[DepSpecKey]
		if !exist {
			msg := fmt.Sprintf("no %s data", DepSpecKey)
			utils.ExitError(log, fmt.Errorf(msg), msg)
		}
		if err := yaml.Unmarshal([]byte(depSpecStr), depSpec); err != nil {
			utils.ExitError(log, err, "cannot unmarshal deployment spec from configmap")
		}
	}

	// Merge DEPLOY_ENV_JSON into deployment
	if depEnvJsonRaw != "" && depEnvJsonRaw != "{}" {
		// Replace quotes
		// (Temporary...as Template CRD does not support JSON string)
		reg, err := regexp.Compile("[']+")
		if err != nil {
			utils.ExitError(log, err, "cannot compile regexp")
		}

		// Unmarshal JSON to env
		depEnvJsonStr := reg.ReplaceAllString(depEnvJsonRaw, "\"")
		depEnv := make(map[string]string)
		if err := json.Unmarshal([]byte(depEnvJsonStr), &depEnv); err != nil {
			utils.ExitError(log, err, "cannot unmarshal DEPLOY_ENV_JSON")
		}

		if len(depSpec.Spec.Template.Spec.Containers) == 0 {
			depSpec.Spec.Template.Spec.Containers = append(depSpec.Spec.Template.Spec.Containers, corev1.Container{})
		}
		mergeEnv(&depSpec.Spec.Template.Spec.Containers[0], depEnv)

		// If CM name is provided, update cm with merged env
		if configMapName != "" {
			depSpecStr, err := yaml.Marshal(depSpec)
			if err != nil {
				utils.ExitError(log, err, "cannot marshal deploy spec")
			}
			cm.Data[DepSpecKey] = string(depSpecStr)
			if err := c.Update(context.TODO(), cm); err != nil {
				utils.ExitError(log, err, "cannot update configMap")
			}
		}
	}

	// Set default values
	setDefaults(depSpec, appName, imageUrl)

	// Add additional env for deployment refresh
	nowTime := time.Now()
	timeEnv := corev1.EnvVar{Name: "DEPLOY_TIME", Value: nowTime.Format("2006-01-02 15:04:05")}
	depSpec.Spec.Template.Spec.Containers[0].Env = append(depSpec.Spec.Template.Spec.Containers[0].Env, timeEnv)

	// Delete ImagePullSecrets if it is empty...
	if len(depSpec.Spec.Template.Spec.ImagePullSecrets) == 1 && depSpec.Spec.Template.Spec.ImagePullSecrets[0].Name == "" {
		depSpec.Spec.Template.Spec.ImagePullSecrets = depSpec.Spec.Template.Spec.ImagePullSecrets[:0]
	}

	// Marshal into YAML
	serializer := k8sjson.NewSerializerWithOptions(k8sjson.DefaultMetaFactory, nil, nil, k8sjson.SerializerOptions{Yaml: true, Pretty: true})
	buf := new(bytes.Buffer)
	if err := serializer.Encode(depSpec, buf); err != nil {
		utils.ExitError(log, err, "cannot marshal deployment into YAML")
	}

	if err := ioutil.WriteFile(TargetFilePath, buf.Bytes(), 0777); err != nil {
		utils.ExitError(log, err, "cannot write file")
	}
	fmt.Println(buf)
}

func mergeEnv(cont *corev1.Container, depEnv map[string]string) {
	for k, v := range depEnv {
		found := false
		for i, e := range cont.Env {
			if e.Name == k {
				cont.Env[i].Value = v
				found = true
				continue
			}
		}
		if !found {
			cont.Env = append(cont.Env, corev1.EnvVar{Name: k, Value: v})
		}
	}
}

func setDefaults(dep *appsv1.Deployment, appName, imageUrl string) {
	// Set Type Meta
	dep.TypeMeta = metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"}

	// Set Object name
	if dep.ObjectMeta.Name == "" {
		dep.ObjectMeta.Name = appName
	}

	// Set labelSelector
	if dep.Spec.Selector == nil {
		dep.Spec.Selector = &metav1.LabelSelector{}
	}
	if dep.Spec.Selector.MatchLabels == nil {
		dep.Spec.Selector.MatchLabels = map[string]string{}
	}
	foundSelector := false
	for k := range dep.Spec.Selector.MatchLabels {
		if k == LabelKey {
			foundSelector = true
		}
	}
	if !foundSelector {
		dep.Spec.Selector.MatchLabels[LabelKey] = appName
	}

	// Set PodTemplate - Label
	if dep.Spec.Template.ObjectMeta.Labels == nil {
		dep.Spec.Template.ObjectMeta.Labels = map[string]string{}
	}
	foundLabel := false
	for k := range dep.Spec.Template.ObjectMeta.Labels {
		if k == LabelKey {
			foundLabel = true
		}
	}
	if !foundLabel {
		dep.Spec.Template.ObjectMeta.Labels[LabelKey] = appName
	}

	// Set Container spec
	if len(dep.Spec.Template.Spec.Containers) == 0 {
		dep.Spec.Template.Spec.Containers = append(dep.Spec.Template.Spec.Containers, corev1.Container{})
	}
	if dep.Spec.Template.Spec.Containers[0].Name == "" {
		dep.Spec.Template.Spec.Containers[0].Name = appName
	}
	if dep.Spec.Template.Spec.Containers[0].Image == "" {
		dep.Spec.Template.Spec.Containers[0].Image = imageUrl
	}
	if dep.Spec.Template.Spec.Containers[0].ImagePullPolicy == "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = corev1.PullAlways
	}
}
