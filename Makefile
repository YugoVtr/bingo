.PHONY: build

default: build

db:
	docker-compose up -d

server:
	go run internal/cmd/server/main.go

build: server
