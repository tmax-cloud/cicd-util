module github.com/cqbqdd11519/cicd-util

go 1.13

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-logr/logr v0.2.0
	github.com/operator-framework/operator-sdk v0.17.1
	github.com/radovskyb/watcher v1.0.7
	github.com/tidwall/gjson v1.6.0
	github.com/tmax-cloud/registry-operator v0.1.2-0.20210119035521-370624778c32
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	sigs.k8s.io/controller-runtime v0.6.4
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.1.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
	k8s.io/api => k8s.io/api v0.18.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.4
)
