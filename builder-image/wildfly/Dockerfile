FROM quay.io/wildfly/wildfly-centos7:18.0

ENV MVN_CENTRAL_URL=https://repo1.maven.org/maven2

USER root
COPY assemble /usr/libexec/s2i/assemble
COPY settings.xml /home/jboss/.m2/settings.xml
RUN chown jboss:root /home/jboss/.m2/settings.xml
USER jboss
