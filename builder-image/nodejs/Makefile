REGISTRY := tmaxcloudck
IMAGE     = s2i-nodejs
VERSION  := 12

URL = $(REGISTRY)/$(IMAGE):$(VERSION)

.PHONY: all build push
all: build push

build:
	docker build -t $(URL) -f ./Dockerfile .

push:
	docker push $(URL)
