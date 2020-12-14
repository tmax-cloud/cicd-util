FROM lachlanevenson/k8s-kubectl:latest

USER root

COPY bin/cicd-util /cicd-util

ENV TZ=Asia/Seoul
RUN apk add tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["/cicd-util"]
