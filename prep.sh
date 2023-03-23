#!/bin/sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"'
rm -f fleek.1
rm -f fleek.1.gz
./fleek man > fleek.1
gzip fleek.1
