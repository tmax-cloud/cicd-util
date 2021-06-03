FROM centos/nodejs-12-centos7

USER root
# Tmax-cloud - install checkfile
RUN curl -L -s -o /usr/bin/checkfile https://github.com/cqbqdd11519/checkfile/releases/download/v0.0.3/checkfile && \
    chmod +x /usr/bin/checkfile
USER 1001
