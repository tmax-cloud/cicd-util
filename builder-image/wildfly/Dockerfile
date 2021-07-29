FROM quay.io/wildfly/wildfly-centos7:18.0

ENV MVN_CENTRAL_URL=https://repo1.maven.org/maven2

USER root
COPY assemble /usr/libexec/s2i/assemble
COPY settings.xml /home/jboss/.m2/settings.xml
RUN chown jboss:root /home/jboss/.m2/settings.xml
# Tmax-cloud - install checkfile
RUN curl -L -s -o /usr/bin/checkfile https://github.com/cqbqdd11519/checkfile/releases/download/v0.0.4/checkfile && \
    chmod +x /usr/bin/checkfile
USER jboss
