REGISTRY := tmaxcloudck
IMAGE     = s2i-django
VERSION  := 35

URL = $(REGISTRY)/$(IMAGE):$(VERSION)

ORIGINAL = centos/python-35-centos7

.PHONY: all build push
all: build push

build:
	docker pull $(ORIGINAL)
	docker tag $(ORIGINAL) $(URL)

push:
	docker push $(URL)
