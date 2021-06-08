REGISTRY := tmaxcloudck
IMAGE     = s2i-apache
VERSION  := 2.4

URL = $(REGISTRY)/$(IMAGE):$(VERSION)

.PHONY: all build push
all: build push

build:
	patch -d ./httpd-container -p0 < tmaxcloud.patch
	docker build --no-cache -t $(URL) -f ./httpd-container/2.4/Dockerfile ./httpd-container/2.4/

push:
	docker push $(URL)
