APP_NAME=go-api-service

run:
	go run ./cmd/api

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

build:
	go build -o bin/api ./cmd/api

docker-build:
	docker build -t $(APP_NAME):local .

docker-up:
	docker compose up --build

docker-down:
	docker compose down