compose-all:
	docker-compose -f ./docker-compose.yml -f ./trades-receiver-service/docker-compose.yml up
compose:
	docker-compose up -d --build
ci:
	golangci-lint run
proto:
	protoc -I api --go_out=api api/*.proto
proto-local:
	protoc -I api -I $${HOME}/github/protobuf/src api/*.proto --go_out=api
