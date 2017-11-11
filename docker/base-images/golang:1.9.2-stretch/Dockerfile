FROM golang:1.9.2-stretch

ENV GOMETALINTER_VERSION="0960299513738ff031fe418b3fcd4f6badc1a095"

COPY scripts/install-dependencies.sh /usr/bin

RUN set -eux; \
	install-dependencies.sh; \
	rm -rf \
		${GOPATH}/src/* \
		${GOPATH}/pkg/*; \
	strip ${GOPATH}/bin/*

RUN set -eux; \
    DEBIAN_FRONTEND=noninteractive; \
    apt-get -q update; \
    apt-get -qy install wget gnupg; \
    wget -qO - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -; \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ stretch-pgdg main" > /etc/apt/sources.list.d/postgres.list; \
    apt-get -q update; \
    apt-get -qy install postgresql-client-9.5; \
    apt-get -qy remove --purge gnupg wget; \
    apt-get -qy autoremove --purge; \
    apt-get -q clean; \
    rm -rf /var/lib/apt/lists/*
