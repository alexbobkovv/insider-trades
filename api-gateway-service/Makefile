compose:
	docker-compose up -d --build
compose-prod:
	docker-compose -f docker-compose.prod.yml up --build
compose-dev:
	docker-compose -f docker-compose.dev.yml up -d --build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v ./cmd/main.go
test:
	go test -v -cover -race ./internal/...
mock:
	go generate -v ./...
ci:
	golangci-lint run
swag:
	swag init -g internal/controller/httpapi/insider_trade.go

.DEFAULT_GOAL := build
