#!/bin/bash
set -e 

go test -race 

go vet ./...

CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -extldflags "static"' -o mimi main.go

docker stop mimi ||true
docker rm mimi ||true

docker run -d --name mimi --restart always \
	--net host \
	--env-file .env \
	-v /etc/localtime:/etc/localtime:ro \
	-v "$PWD"/mimi:/mimi \
	-w / \
	debian:jessie ./mimi

docker logs -f mimi
