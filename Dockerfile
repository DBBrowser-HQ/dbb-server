FROM golang:latest

ENV GOPATH=/
COPY ./ ./

RUN go mod download
RUN make

CMD ["./dbb-server"]
