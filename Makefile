all: build

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_TIME := $(shell git show -s --format=%ci $(GIT_COMMIT) | tr -d '\n')
GIT_TAG := $(shell git describe --exact-match --tags $(GIT_COMMIT) 2> /dev/null)

PACKAGE := github.com/FallenTaters/streepjes/src/infrastructure/router

LDFLAGS := '-X "$(PACKAGE).buildCommit=$(GIT_COMMIT)" \
		-X "$(PACKAGE).buildTime=$(GIT_TIME)" \
		-X "$(PACKAGE).buildVersion=$(GIT_TAG)"'

generate:
	@echo "Generating code..."
	@go generate ./...
	@echo "Done"

vugugen:
	@echo "Running vugugen..."
	@bash ./frontend/generate.bash
	@echo "Done"

wasm:
	@echo "Compiling frontend..."
	@GOARCH=wasm GOOS=js go build -o ./static/files/app.wasm ./frontend/
	@echo "Done"

run: vugugen wasm
	@echo "Starting local server..."
	@go run -ldflags "-X $(PACKAGE).buildVersion=development" .

run-backend:
	@go run -ldflags "-X $(PACKAGE).buildVersion=development" .

test:
	@echo "Testing..."
	@go test ./backend/... -cover

lint:
	@echo "Linting..."
	@golangci-lint run ./...
	@echo "Done"

build: generate vugugen wasm
	@go build -o ./bin/streepjes -ldflags $(LDFLAGS) .
