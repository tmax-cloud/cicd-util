REGISTRY := tmaxcloudck
IMAGE     = s2i-tomcat
VERSION  := 8.5

URL = $(REGISTRY)/$(IMAGE):$(VERSION)

SUBMODULE_PATH = ./s2i-tomcat/tomcat-8.5/maven-3.5.0/jdk-8

.PHONY: all build push
all: build push

build:
	patch -d ./s2i-tomcat -p0 < tmaxcloud.patch
	chmod +x $(SUBMODULE_PATH)/s2i/bin/*
	docker build --no-cache -t $(URL) -f $(SUBMODULE_PATH)/Dockerfile $(SUBMODULE_PATH)

push:
	docker push $(URL)
