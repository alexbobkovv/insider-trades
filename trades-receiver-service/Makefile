compose:
	docker-compose up -d --build
compose-prod:
	docker-compose -f docker-compose.prod.yml up --build
compose-dev:
	docker-compose -f docker-compose.dev.yml up -d --build
build:
	go build -v ./cmd/main.go
test:
	go test -v -cover -race ./internal/...
mock:
	go generate -v ./...
ci:
	golangci-lint run
migrate-up:
	docker-compose --profile tools run migrate up
migrate-down:
	docker-compose --profile tools run migrate down
swag:
	swag init -g internal/controller/httpapi/insider_trade.go

.DEFAULT_GOAL := build

