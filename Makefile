compose-all:
	docker-compose -f ./docker-compose.yml -f ./trades-receiver-service/docker-compose.yml up
