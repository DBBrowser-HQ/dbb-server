ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: build
build:
	GOCACHE=`pwd`/.cache go build -v -o dbb-server ./cmd/dbb

.PHONY: dRun
dRun:
	docker-compose up --build -d dbb-server

.PHONY: mUp
mUp:
	migrate -path ./migrations -database \
 'postgres://$(SERVER_DB_USERNAME):$(SERVER_DB_PASSWORD)@$(SERVER_DB_HOST):$(SERVER_DB_PORT)/$(SERVER_DB_NAME)?sslmode=disable' up

.PHONY: mDown
mDown:
	migrate -path ./migrations -database \
 'postgres://$(SERVER_DB_USERNAME):$(SERVER_DB_PASSWORD)@$(SERVER_DB_HOST):$(SERVER_DB_PORT)/$(SERVER_DB_NAME)?sslmode=disable' down

.DEFAULT_GOAL := build