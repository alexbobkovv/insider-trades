compose-all:
	docker-compose -f ./docker-compose.yml -f ./trades-receiver-service/docker-compose.yml up
compose:
	docker-compose -f docker-compose.yml up --build
test:
	go test -v -cover -race ./...
ci:
	golangci-lint run
proto:
	protoc -I api --go_out=api api/*.proto --go-grpc_out=api
proto-local:
	protoc -I api -I $${HOME}/github/protobuf/src api/*.proto --go_out=api --go-grpc_out=api
build-receiver:
	docker build --build-arg SERVICE_NAME=trades-receiver-service -t trades-receiver-service .
run-receiver:
	docker run -d -p 8080:8080 -p 50051:50051 --name trades-receiver-service --network=insider-trades --env-file=trades-receiver-service/.env trades-receiver-service
