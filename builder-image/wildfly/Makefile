REGISTRY := tmaxcloudck
IMAGE     = s2i-wildfly
VERSION  := 18

URL = $(REGISTRY)/$(IMAGE):$(VERSION)

.PHONY: all build push
all: build push

build:
	docker build --no-cache -t $(URL) -f ./Dockerfile .

push:
	docker push $(URL)
