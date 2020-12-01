build:
	go build -a -installsuffix cgo -o zero-notification-service ./cmd/server/

run:
	go run cmd/server/main.go

generate:
	openapi-generator generate -i api/notification-service.yaml -g go-server -o ./ -p sourceFolder=internal/server -p packageName=server --git-user-id=commitdev --git-repo-id=zero-notification-service
	go get golang.org/x/tools/cmd/goimports
	goimports -w internal/server/

docker:
	docker build -t zero-notification-service .

docker-run:
	docker run --rm -it zero-notification-service
