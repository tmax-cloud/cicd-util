FROM tmaxsoftofficial/jeus:8.1.105067.1-jdk8-openjdk-ubuntu

LABEL io.openshift.s2i.scripts-url="image:///usr/libexec/s2i"

ENV MVN_CENTRAL_URL https://repo1.maven.org/maven2
ENV HOME /opt/app-root/src

RUN apt update && \
    apt install zip maven -y && \
    apt-get clean autoclean && \
    apt-get autoremove --yes && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /usr/libexec && \
    mkdir -p /opt/s2i/destination && \
    mkdir -p /opt/app-root/src/.m2 && \
    mkdir -p /opt/webapps && \
    mkdir -p /opt/destination && \
    rm -rf /root/app/examples.ear && \
    ln -s /opt/destination/examples.ear /root/app/examples.ear

COPY s2i-scripts /usr/libexec/s2i
COPY container-entrypoint /usr/bin/container-entrypoint
COPY settings.xml /opt/app-root/src/.m2/settings.xml
COPY examples.ear /opt/destination/examples.ear

RUN chmod +x /usr/libexec/s2i/* && \
    chmod +x /usr/bin/container-entrypoint && \
    chown -R 1001:1001 /usr/libexec/s2i/ && \
    chown -R 1001:1001 /usr/bin/container-entrypoint && \
    chown -R 1001:1001 /opt/app-root && \
    chown -R 1001:1001 /opt/s2i && \
    chown -R 1001:1001 /opt/webapps && \
    chown -R 1001:1001 /opt/destination && \
    chown -R 1001:1001 /root

WORKDIR /opt/app-root/src

EXPOSE 8080

ENTRYPOINT ["/usr/bin/container-entrypoint"]
