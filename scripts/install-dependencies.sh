#!/bin/bash -eux

GOMETALINTER_VERSION="0960299513738ff031fe418b3fcd4f6badc1a095"

# dep
go get github.com/golang/dep/cmd/dep

# goveralls
go get github.com/mattn/goveralls


# gometalinter
test ! -d ${GOPATH}/src/github.com/alecthomas/gometalinter && \
	git clone https://github.com/alecthomas/gometalinter.git ${GOPATH}/src/github.com/alecthomas/gometalinter
cd ${GOPATH}/src/github.com/alecthomas/gometalinter
git checkout -B build
git reset --hard ${GOMETALINTER_VERSION}
cd ${GOPATH}
go install github.com/alecthomas/gometalinter
gometalinter --debug --install
