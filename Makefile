PROJECT := github.com/MaxFando/rate-limiter
GIT_COMMIT := $(shell git rev-parse HEAD)
LOCAL_BIN := $(shell pwd)/bin

appName = github.com/MaxFando/rate-limiter
compose = docker-compose -p rate-limiter

install-deps:
	GOBIN=$(LOCAL_BIN) go install go install go.uber.org/mock/mockgen

structurizer:
	docker-compose -f docs/structurizer/docker-compose.yml up -d

run: up

up: down build
	@echo "Starting app..."
	$(compose) up -d
	@echo "Docker images built and started!"

build:
	@echo "Building images"
	@echo "Docker images built!"

down:
	@echo "Stopping docker compose..."
	$(compose) down -v
	@echo "Done!"

test:
	go test -race -count 100 -v ./internal/service/...

lint:
	golangci-lint run -c .golangci.yaml

lint-fix:
	golangci-lint run -v -c .golangci.yaml --fix ./...

generate:
	@echo "Generating code..."
	go generate ./...
	@echo "Code generated!"

migrate:
	migrate -version $(version)

migrate.down:
	migrate -source file://migrations -database postgres://localhost:5432/postgres?sslmode=disable down

migrate.up:
	migrate -source file://migrations -database postgres://localhost:5432/postgres?sslmode=disable up

proto:
	protoc --proto_path=internal/delivery/grpcapi/proto/blacklist internal/delivery/grpcapi/proto/blacklist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/delivery/grpcapi/proto/whitelist internal/delivery/grpcapi/proto/whitelist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/delivery/grpcapi/proto/bucket internal/delivery/grpcapi/proto/bucket/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/delivery/grpcapi/proto/auth internal/delivery/grpcapi/proto/auth/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import

deploy:
	@(echo "Deploying to k8s...")
	@(echo "Deployed!")