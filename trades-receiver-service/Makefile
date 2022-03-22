build:
	go build -v ./cmd/main.go
test:
	go test -v -cover -race ./internal/...
ci:
	golangci-lint run
mock:
	go generate -v ./...
.DEFAULT_GOAL := build

