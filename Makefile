.PHONY: build

default: build

dev:
	go run cmd/web/server.go

db:
	docker compose up -d db

server:
	docker compose up -d server

test:
	go test ./... -count=1 -cover -timeout=30s | grep -v ?

test_integration:
	go test ./... -timeout 100s -count=5 --tags=integration | grep -v ?

lint:
	golangci-lint run --timeout=5m

ci: test lint

build: server
