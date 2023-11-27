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

build-dev:
	go build -o $(BIN_APP) $(CMD_APP_PATH)/main.go

run-dev: build
	$(BIN_APP) -config ./configs/config.yaml

build:
	docker-compose -f ${DOCKER_COMPOSE} up -d --build

run:
	docker-compose -f ${DOCKER_COMPOSE} up

down:
	docker-compose -f ${DOCKER_COMPOSE} down

build-admin:
	go build -o $(BIN_ADMIN) $(CMD_ADMIN_PATH)/main.go

.PHONY: generate test lint build build-dev build-admin runbuild-dev run-admin docker-build docker-down