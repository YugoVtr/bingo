.PHONY: build

default: build

dev:
	go run cmd/web/server.go

db:
	docker compose up -d db

server:
	docker compose up -d server

test:
	docker compose run --rm test

lint:
	docker compose run --rm lint

ci: test lint

build: server
