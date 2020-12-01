FROM openapitools/openapi-generator-cli:v5.0.0-beta3 AS generate
WORKDIR /go/src
copy . .

RUN docker-entrypoint.sh generate -i api/notification-service.yaml -g go-server -o ./ -p sourceFolder=internal/server -p packageName=server --git-user-id=commitdev --git-repo-id=zero-notification-service


FROM golang:1.14 AS build
WORKDIR /go/src

COPY go.mod .
COPY go.sum .

## fetch dependencies
RUN go mod tidy && go mod download

## copy project files
COPY . .
COPY --from=generate /go/src/ ./

RUN go get golang.org/x/tools/cmd/goimports && goimports -w internal/server/

ENV CGO_ENABLED=0


RUN make build


FROM scratch AS runtime
COPY --from=build /go/src/zero-notification-service ./
EXPOSE 80/tcp
ENTRYPOINT ["./zero-notification-service"]
