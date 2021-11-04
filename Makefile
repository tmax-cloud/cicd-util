PACKAGE_NAME = github.com/tmax-cloud/cicd-util

REGISTRY ?= tmaxcloudck
VERSION  ?= latest

CICD_UTIL_IMAGE_NAME = cicd-util
CICD_UTIL_IMAGE      = $(REGISTRY)/$(CICD_UTIL_IMAGE_NAME):$(VERSION)

.PHONY: all
all: build image push

cicd-util: bin/cicd-util image-cicd-util push-cicd-util

.PHONY: build
build: bin/cicd-util

bin/%: cmd/%
	CGO_ENABLED=0 go build -o $@ $(PACKAGE_NAME)/$<


.PHONY: image image-cicd-util
image: image-cicd-util

image-cicd-util:
	docker build -f build/cicd-util/Dockerfile -t $(CICD_UTIL_IMAGE) .


.PHONY: tag-latest tag-latest-cicd-util
tag-latest: tag-latest-cicd-util

tag-latest-cicd-util:
	docker tag $(CICD_UTIL_IMAGE) $(REGISTRY)/$(CICD_UTIL_IMAGE_NAME):latest


.PHONY: push push-cicd-util
push: push-cicd-util

push-cicd-util:
	docker push $(CICD_UTIL_IMAGE)


.PHONY: push-latest push-latest-cicd-util
push-latest: push-latest-cicd-util

push-latest-cicd-util:
	docker push $(REGISTRY)/$(CICD_UTIL_IMAGE_NAME):latest


.PHONY: test test-verify save-sha-mod compare-sha-mod verify test-unit test-lint
test: test-verify test-unit test-lint

test-verify: save-sha-mod verify compare-sha-mod

save-sha-mod:
	$(eval MODSHA=$(shell sha512sum go.mod))
	$(eval SUMSHA=$(shell sha512sum go.sum))

verify:
	go mod verify

compare-sha-mod:
	$(eval MODSHA_AFTER=$(shell sha512sum go.mod))
	$(eval SUMSHA_AFTER=$(shell sha512sum go.sum))
	@if [ "${MODSHA_AFTER}" = "${MODSHA}" ]; then echo "go.mod is not changed"; else echo "go.mod file is changed"; exit 1; fi
	@if [ "${SUMSHA_AFTER}" = "${SUMSHA}" ]; then echo "go.sum is not changed"; else echo "go.sum file is changed"; exit 1; fi

test-unit:
	go test -v ./pkg/...

test-lint:
	golangci-lint run ./... -v -E gofmt --timeout 1h0m0s


.PHONY: builder-images
builder-images:
	make -C builder-image/apache
	make -C builder-image/django
	make -C builder-image/nodejs
	make -C builder-image/tomcat
	make -C builder-image/wildfly
