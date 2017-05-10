#!/bin/bash
set -e 

go test -race 

go vet ./...

CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -extldflags "static"' -o whoami main.go

docker stop whoami ||true
docker rm whoami ||true

docker run -d --name whoami --restart always \
	--net host \
	--env-file .env \
	-v /etc/localtime:/etc/localtime:ro \
	-v "$PWD"/whoami:/whoami \
	-w / \
	debian:jessie ./whoami

docker logs -f whoami
