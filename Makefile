.PHONY: all test-usecase test-service

all: usecase service

usecase:
	go test ./internal/usecase -v

service:
	go test ./internal/service -v
