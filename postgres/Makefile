.PHONY: build run stop clean logs
.SILENT: build run stop clean logs

init:
	docker network create finboxnet

build:
	docker build -t postgresql .

run:
	docker run --name postgres -e POSTGRES_PASSWORD=postgres --network finboxnet -p 5432:5432 --hostname localhost -v postgres-volume:/data -d postgresql

stop:
	docker stop $(shell docker ps | grep "postgres" | awk '{ print $$1 }')

clean:
	docker container prune
	kill -9 $(shell lsof -t -i:5432)

logs:
	docker logs $(shell docker ps | grep 'postgres' | awk '{ print $$1 }') -f