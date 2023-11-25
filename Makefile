GO_GEN_PATH="./internal/transport/grpc"
CMD_APP_PATH="cmd/app"
CMD_ADMIN_PATH="cmd/admin-cli"
BIN_APP="bin/app"
BIN_ADMIN="bin/admin-cli"
DOCKER_COMPOSE="build/local/docker-compose.yml"

generate:
	protoc -I. ./docs/proto/api.proto \
		--go_out="$(GO_GEN_PATH)" \
		--go-grpc_out="$(GO_GEN_PATH)"

test:
	go test ./... -v -race -count 100

lint:
	golangci-lint run ./...

build:
	go build -o $(BIN_APP) $(CMD_APP_PATH)/main.go

run: build
	$(BIN_APP) -config ./configs/config.yaml

build-admin:
	go build -o $(BIN_ADMIN) $(CMD_ADMIN_PATH)/main.go

docker-build:
	docker-compose -f ${DOCKER_COMPOSE} up -d --build

docker-up: docker-build
	docker-compose -f ${DOCKER_COMPOSE} up

docker-down:
	docker-compose -f ${DOCKER_COMPOSE} down

.PHONY: generate test lint build build-admin run run-admin docker-build docker-up docker-down