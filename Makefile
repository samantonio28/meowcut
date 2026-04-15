.PHONY: all usecase service lint

all: usecase service

usecase:
	go test ./internal/usecase -v

service:
	go test ./internal/service -v

lint:
	golangci-lint run ./...
