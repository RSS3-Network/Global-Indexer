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
    		--tag rss3-network/global-indexer:$(VERSION) \
    		.

run:
	ENVIRONMENT=development go run ./cmd
