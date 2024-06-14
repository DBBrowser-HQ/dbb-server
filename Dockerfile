FROM golang:1.22.3-alpine3.20

RUN apk update && apk add --update docker openrc && apk add make && apk add curl && apk add --no-cache --upgrade bash
RUN rc-update add docker boot

ENV GOPATH=/
RUN go env -w GOCACHE=/.cache

COPY ./ ./

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go mod download
RUN --mount=type=cache,target=/.cache go build -v -o dbb-server ./cmd/dbb

ENTRYPOINT docker build -f ./Dockerfile_Postgres -t postgres-image . && ./dbb-server
