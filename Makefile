SCRIPT_FOLDER_NAME = commands

run:
	docker compose build
	docker compose up

build:
	docker build -t olx-clone . && docker run -p 8080:8080 olx-clone

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

gen:
	cd $(SCRIPT_FOLDER_NAME) && \
	go run *.go $n
	cd ..

deploy: 
	echo "TODO"

test: 
	echo "TODO"

.PHONY: build run logs dockerstop
.SILENT: build run logs dockerstop