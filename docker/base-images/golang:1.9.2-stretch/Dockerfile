FROM golang:1.9.2-stretch

ENV GOMETALINTER_VERSION="0960299513738ff031fe418b3fcd4f6badc1a095"

RUN set -eux; \
	go get github.com/golang/dep/cmd/dep; \
	\
	git clone https://github.com/alecthomas/gometalinter.git ${GOPATH}/src/github.com/alecthomas/gometalinter; \
	cd ${GOPATH}/src/github.com/alecthomas/gometalinter; \
	git reset --hard ${GOMETALINTER_VERSION}; \
	cd ${GOPATH}; \
	go install github.com/alecthomas/gometalinter; \
	gometalinter --debug --install; \
	\
	rm -rf \
		${GOPATH}/src/* \
		${GOPATH}/pkg/*; \
	strip ${GOPATH}/bin/*
