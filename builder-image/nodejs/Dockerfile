FROM centos/nodejs-12-centos7

USER root
# Tmax-cloud - install checkfile
RUN curl -L -s -o /usr/bin/checkfile https://github.com/tmax-cloud/checkfile/releases/download/v0.0.4/checkfile && \
    chmod +x /usr/bin/checkfile
USER 1001
