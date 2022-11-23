REGISTRY := tmaxcloudck
IMAGE     = nexus-hyperauth
VERSION  := v1.0.0

URL = $(REGISTRY)/$(IMAGE):$(VERSION)

.PHONY: all build push
all: build push

build:
	docker build --no-cache -t $(URL) -f ./Dockerfile .

push:
	docker push $(URL)
