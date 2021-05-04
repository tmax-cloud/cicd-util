FROM centos/python-35-centos7

USER root
# Tmax-cloud - install checkfile
RUN curl -L -s -o /usr/bin/checkfile https://github.com/cqbqdd11519/checkfile/releases/download/v0.0.2/checkfile && \
    chmod +x /usr/bin/checkfile
USER 1001
