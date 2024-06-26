ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: build
build:
	GOCACHE=`pwd`/.cache go build -v -o dbb-server ./cmd/dbb

.PHONY: dockerUp
dockerUp:
	docker compose up --build -d

.PHONY: dockerDown
dockerDown:
	docker compose down -v -t 30

.PHONY: dockerStop
dockerStop:
	docker compose stop

.PHONY: dockerStart
dockerStart:
	docker compose start

ifeq (migration,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: migration
migration:
	migrate -path ./migrations -database 'postgres://$(SERVER_DB_USERNAME):$(SERVER_DB_PASSWORD)@$(SERVER_DB_HOST):$(SERVER_DB_PORT)/$(SERVER_DB_NAME)?sslmode=$(SERVER_DB_SSL_MODE)' $(RUN_ARGS)

.PHONY: migrationUpDownUp
migrationUpDownUp:
	make migration up 1; make migration down 1; make migration up 1

.DEFAULT_GOAL := build