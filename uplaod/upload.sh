#!/bin/sh
set -euv

test -d initial || mkdir -p initial

cp ../cmd/center_config.json initial/center_config.json
#cp ../cmd/executor_config.json initial/executor_config.json

cd initial
GOOS=linux GOARCH=amd64 go build ../../cmd/center.go
GOOS=linux GOARCH=amd64 go build ../../cmd/executor.go
cd ../

tar -zcvf initial.tar.gz initial
scp -r -P 59672 initial.tar.gz common@212.129.131.27:~
ssh -p 59672 common@212.129.131.27 "tar -zxvf initial.tar.gz;rm initial.tar.gz"
rm initial.tar.gz