.PHONY: build

default: build

db:
	docker compose up -d

server:
	go run cmd/web/server.go

build: server
