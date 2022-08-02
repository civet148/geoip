#SHELL=/usr/bin/env bash

CLEAN:=
BINS:=
DATE_TIME=`date +'%Y%m%d %H:%M:%S'`
COMMIT_ID=`git rev-parse --short HEAD`
CODE_DIR=${PWD}
BIN_NAME=geoip

build:
	rm -f ${BIN_NAME}
	git submodule update --init --recursive \
	&& go build -ldflags "-s -w -X 'main.BuildTime=${DATE_TIME}' -X 'main.GitCommit=${COMMIT_ID}'" -o ${BIN_NAME} cmd/main.go
.PHONY: build
BINS+=${BIN_NAME}

ip-data:
	tar xvfz ip.tar.gz

test:build
	./geoip test
.PHONY: test

import:build ip-data
	./geoip import --dsn "root:123456@tcp(192.168.20.108:3306)/ip-system?charset=utf8"  --table ip_info --data ip.dat

clean:
	rm -rf $(CLEAN) $(BINS)
.PHONY: clean
