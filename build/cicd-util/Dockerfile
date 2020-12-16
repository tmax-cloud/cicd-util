FROM golang:1.14 AS s2i
ENV TAG=v1.3.1
COPY s2i.patch /tmp/s2i.patch
RUN apt update && \
    apt install -y patch && \
    mkdir -p /tmp/s2i && \
    cd /tmp/s2i && \
    git clone --branch $TAG --depth 1 https://github.com/openshift/source-to-image.git . && \
    patch -p0 < /tmp/s2i.patch && \
    make

FROM lachlanevenson/k8s-kubectl:latest

USER root

COPY --from=s2i /tmp/s2i/_output/local/bin/linux/amd64/s2i /usr/local/bin/s2i
COPY bin/cicd-util /cicd-util

ENV TZ=Asia/Seoul
RUN apk add tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["/cicd-util"]
