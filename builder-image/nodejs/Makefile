REGISTRY := tmaxcloudck
IMAGE     = s2i-nodejs
VERSION  := 12

URL = $(REGISTRY)/$(IMAGE):$(VERSION)

ORIGINAL = centos/nodejs-12-centos7

.PHONY: all build push
all: build push

build:
	docker pull $(ORIGINAL)
	docker tag $(ORIGINAL) $(URL)

push:
	docker push $(URL)
