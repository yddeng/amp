#!/bin/sh
set -euv

test -d initial || mkdir -p initial

cp ../server/config.json initial/config.json
cp ../server/nav.json initial/nav.json

cd initial
GOOS=linux GOARCH=amd64 go build ../../server/server.go
cd ../

tar -zcvf initial.tar.gz initial
scp -r -P 59672 initial.tar.gz common@212.129.131.27:~
ssh -p 59672 common@212.129.131.27 "tar -zxvf initial.tar.gz;rm initial.tar.gz"
rm initial.tar.gz