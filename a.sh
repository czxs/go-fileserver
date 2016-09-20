#!/bin/sh
export GOPATH="/home/rabbitmq/fileserver/"
cd /home/rabbitmq/fileserver/src/
go build
rsync -azv src /usr/local/fileserver/src/
rsync -azv src /usr/local/fileserver1/src/

cd -
