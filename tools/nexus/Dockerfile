# https://github.com/sonatype/docker-nexus3
# https://hub.docker.com/r/sonatype/nexus3
FROM sonatype/nexus3:3.22.0
USER root

ENV NEXUS_HOME="/opt/sonatype/nexus"
ENV NEXUS_PLUGINS ${NEXUS_HOME}/system

WORKDIR "${NEXUS_HOME}"

# https://github.com/flytreeleft/nexus3-keycloak-plugin
ENV KEYCLOAK_PLUGIN_VERSION 0.6.0-SNAPSHOT
ENV KEYCLOAK_PLUGIN_RELEASE 0.6.0-SNAPSHOT
ENV KEYCLOAK_PLUGIN /org/github/flytreeleft/nexus3-keycloak-plugin/${KEYCLOAK_PLUGIN_RELEASE}


ADD ./nexus3-keycloak-plugin-${KEYCLOAK_PLUGIN_VERSION}.jar \
    ${NEXUS_PLUGINS}/org/github/flytreeleft/nexus3-keycloak-plugin/${KEYCLOAK_PLUGIN_VERSION}/nexus3-keycloak-plugin-${KEYCLOAK_PLUGIN_VERSION}.jar

RUN chmod 644 ${NEXUS_PLUGINS}/org/github/flytreeleft/nexus3-keycloak-plugin/${KEYCLOAK_PLUGIN_VERSION}/nexus3-keycloak-plugin-${KEYCLOAK_PLUGIN_VERSION}.jar
RUN chown nexus ${NEXUS_PLUGINS}/org/github/flytreeleft/nexus3-keycloak-plugin/${KEYCLOAK_PLUGIN_VERSION}/nexus3-keycloak-plugin-${KEYCLOAK_PLUGIN_VERSION}.jar
RUN chown nexus:nexus -R /opt/sonatype/nexus

EXPOSE 5000 8081 8443
USER nexus

CMD ["bin/nexus", "run"]
