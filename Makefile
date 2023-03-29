VERSION ?= SNAPSHOT
BUILD ?=$(shell git rev-parse --short HEAD)
PKG ?=github.com/commitdev/zero-notification-service
BUILD_ARGS=-installsuffix cgo -v -trimpath -ldflags=all="-X main.appVersion=${VERSION} -X main.appBuild=${BUILD}"

build:
	go build ${BUILD_ARGS} -o zero-notification-service ./cmd/server/

run:
	go run cmd/server/main.go

generate:
	openapi-generator generate -i api/notification-service.yaml -g go-server -o ./ -p sourceFolder=internal/server -p packageName=server --git-user-id=commitdev --git-repo-id=zero-notification-service
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w internal/server/

docker:
	docker build -t zero-notification-service .

docker-run:
	docker run --rm -it zero-notification-service

test:
	go test ./...
