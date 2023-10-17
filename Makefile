run:
	docker build -t olx-clone . && docker run -p 8080:8080 olx-clone

build:
	docker build -t olx-clone

start:
	docker run -p 8080:8080 olx-clone

compose:
	docker compose build && docker compose up

down:
	docker compose down

dev:
	nodemon --exec go run main.go

install:
	go mod tidy

deploy: 
	echo "TODO"

test: 
	echo "TODO"

.PHONY: build run logs dockerstop
.SILENT: build run logs dockerstop