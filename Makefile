PROJECT := github.com/MaxFando/rate-limiter
GIT_COMMIT := $(shell git rev-parse HEAD)

appName = github.com/MaxFando/rate-limiter
compose = docker-compose -p rate-limiter

structurizer:
	docker-compose -f docs/structurizer/docker-compose.yml up -d

run: up

up: down build
	@echo "Starting app..."
	$(compose) up -d
	@echo "Docker images built and started!"

build:
	@echo "Building images"
	$(compose) build
	@echo "Docker images built!"

down:
	@echo "Stopping docker compose..."
	$(compose) down
	@echo "Done!"

test:
	go test -race ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm 	coverage.out

lint:
	golangci-lint run -c .golangci.yaml

mock:
	@echo "Generating mocks..."
	rm -rf internal/mocks
	mockery --all --case unserscore --keeptree --dir internal/service --output mocks/service --log-level warn
	mockery --all --case unserscore --keeptree --dir internal/usecase --output mocks/usecase --log-level warn
	@echo "Mocks generated!"

migrate:
	migrate -version $(version)

migrate.down:
	migrate -source file://migrations -database postgres://localhost:5432/postgres?sslmode=disable down

migrate.up:
	migrate -source file://migrations -database postgres://localhost:5432/postgres?sslmode=disable up

proto:
	protoc --proto_path=internal/delivery/grpc/proto/blacklist internal/delivery/grpc/proto/blacklist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/delivery/grpc/proto/whitelist internal/delivery/grpc/proto/whitelist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/delivery/grpc/proto/bucket internal/delivery/grpc/proto/bucket/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/delivery/grpc/proto/auth internal/delivery/grpc/proto/auth/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import

