BIN := ./bin/calendar
DOCKER_IMG="calendar:develop"
PROTO_VERSION := v1
PROTO_OUTPUT_DIR := pkg/calendar-api
PROTO_MODULE_NAME := github.com/devv4n/otus-hw/hw12_13_14_15_calendar/pkg/calendarapi

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN)_scheduler -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN)_sender -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config ./configs/config.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.0

lint: install-lint-deps
	golangci-lint run ./... -v

lint-fix:
	gofumpt -l -w .
	golangci-lint run ./... --fix

generate: ## Generate pb files from given proto files
	mkdir -p $(PROTO_OUTPUT_DIR)
	rm -f $(PROTO_OUTPUT_DIR)/*.go
	protoc \
		-I=api \
		-I=third_party/googleapis \
		--go_out=$(PROTO_OUTPUT_DIR) \
		--go_opt=module=$(PROTO_MODULE_NAME) \
		--go-grpc_out=$(PROTO_OUTPUT_DIR) \
		--go-grpc_opt=module=$(PROTO_MODULE_NAME) \
		--grpc-gateway_out=$(PROTO_OUTPUT_DIR) \
		--grpc-gateway_opt=module=$(PROTO_MODULE_NAME) \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=generate_unbound_methods=true \
		./api/EventService.proto

up:
	docker-compose -f ./deployments/docker-compose.yaml up -d --build

down:
	docker-compose -f ./deployments/docker-compose.yaml  down

integration-tests:
	@echo "Запуск интеграционных тестов..."
	@set -e; \
	trap 'docker-compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.test.yaml down --volumes --remove-orphans' EXIT; \
	docker-compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.test.yaml up --build --abort-on-container-exit --exit-code-from integration_tests
	@echo "Интеграционные тесты завершены."

.PHONY: build run build-img run-img version test lint generate lint-fix up down integration-tests
