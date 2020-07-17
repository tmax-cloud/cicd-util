PACKAGE_NAME = github.com/cqbqdd11519/cicd-util

REGISTRY := tmaxcloudck
VERSION  := 1.1.0

CICD_UTIL_IMAGE_NAME = cicd-util
CICD_UTIL_IMAGE      = $(REGISTRY)/$(CICD_UTIL_IMAGE_NAME):$(VERSION)

SONAR_CLIENT_IMAGE_NAME = sonar-client
SONAR_CLIENT_IMAGE      = $(REGISTRY)/$(SONAR_CLIENT_IMAGE_NAME):$(VERSION)


.PHONY: all
all: build image push

cicd-util: bin/cicd-util image-cicd-util push-cicd-util

sonar-client: bin/sonar-client image-sonar-client push-sonar-client

.PHONY: build
build: bin/cicd-util bin/sonar-client

bin/%: cmd/%
	CGO_ENABLED=0 go build -o $@ $(PACKAGE_NAME)/$<


.PHONY: image image-cicd-util image-sonar-client
image: image-cicd-util image-sonar-client

image-cicd-util:
	docker build -f build/cicd-util/Dockerfile -t $(CICD_UTIL_IMAGE) .

image-sonar-client:
	docker build -f build/sonar-client/Dockerfile -t $(SONAR_CLIENT_IMAGE) .


.PHONY: push push-cicd-util push-sonar-client
push: push-cicd-util push-sonar-client

push-cicd-util:
	docker push $(CICD_UTIL_IMAGE)

push-sonar-client:
	docker push $(SONAR_CLIENT_IMAGE)
