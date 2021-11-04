module github.com/tmax-cloud/cicd-util

go 1.13

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-logr/logr v0.1.0
	github.com/operator-framework/operator-sdk v0.17.1
	github.com/radovskyb/watcher v1.0.7
	github.com/tidwall/gjson v1.6.0
	k8s.io/api v0.17.6
	k8s.io/apimachinery v0.17.6
	sigs.k8s.io/controller-runtime v0.5.2
)

replace k8s.io/client-go => k8s.io/client-go v0.17.4
