.PHONY: build

default: build

db:
	docker compose up -d

server:
	go run cmd/web/server.go

test:
	go test ./... -count=1 -cover | grep -v ?

build: server
