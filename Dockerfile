FROM golang:latest

ENV GOPATH=/
RUN go env -w GOCACHE=/.cache

COPY ./ ./

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go mod download
RUN --mount=type=cache,target=/.cache go build -v -o dbb-server ./cmd/dbb

CMD ["./dbb-server"]
