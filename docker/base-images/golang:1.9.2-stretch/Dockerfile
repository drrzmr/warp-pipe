FROM golang:1.9.2-stretch

ENV GOMETALINTER_VERSION="0960299513738ff031fe418b3fcd4f6badc1a095"

COPY scripts/install-dependencies.sh /usr/bin

RUN set -eux; \
	install-dependencies.sh; \
	rm -rf \
		${GOPATH}/src/* \
		${GOPATH}/pkg/*; \
	strip ${GOPATH}/bin/*
