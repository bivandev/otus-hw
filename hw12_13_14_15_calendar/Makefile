BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
PROTO_VERSION := v1
PROTO_OUTPUT_DIR := pkg/calendar-api
PROTO_MODULE_NAME := github.com/devv4n/otus-hw/hw12_13_14_15_calendar/pkg/calendarapi

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

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
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./... -v --concurrency 1

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

.PHONY: build run build-img run-img version test lint generate lint-fix
