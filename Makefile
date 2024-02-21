VERSION=$(shell git describe --tags --abbrev=0)

ifeq ($(VERSION),)
	VERSION="0.0.0"
endif

lint:
	go mod tidy
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run

test:
	go test -cover -race -v ./...

.PHONY: build
build:
	mkdir -p ./build
	go build \
		-o ./build/rss3-global-indexer ./cmd

image:
	docker build \
    		--tag naturalselectionlabs/rss3-global-indexer:$(VERSION) \
    		.

run:
	ENVIRONMENT=development go run ./cmd

OAPI_SPEC ?= spec/gateway.yml
OAPI_TARGET ?= internal/service/gateway/gen/oapi/
OAPI_TARGET_FILENAME ?= oapi.go
gengatewayapi:
	mkdir -p $(OAPI_TARGET)
	go get github.com/deepmap/oapi-codegen/v2
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen -package oapi -generate=types,client,server,spec,skip-prune -o $(OAPI_TARGET)$(OAPI_TARGET_FILENAME) $(OAPI_SPEC)
	go mod tidy
