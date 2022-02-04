all: build

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_TIME := $(shell git show -s --format=%ci $(GIT_COMMIT) | tr -d '\n')
GIT_TAG := $(shell git describe --exact-match --tags $(GIT_COMMIT) 2> /dev/null)

PACKAGE := github.com/PotatoesFall/vecty-test/src/infrastructure/router

LDFLAGS := '-X "$(PACKAGE).buildCommit=$(GIT_COMMIT)" \
		-X "$(PACKAGE).buildTime=$(GIT_TIME)" \
		-X "$(PACKAGE).buildVersion=$(GIT_TAG)"'

generate:
	@go generate ./...

wasm:
	@GOARCH=wasm GOOS=js go build -o ./static/files/app.wasm ./frontend/

run: generate wasm
	@go run -ldflags "-X $(PACKAGE).buildVersion=development" .

arelo:
	@arelo -p '**/*.go' -p '**/*.css' -p '**/*.js' -p '**/*.html' -i '**/.*' -i '**/*_test.go' -i '**/*_enumer.go' -- make run

build: generate wasm
	@go build -o ./bin/streepjes -ldflags $(LDFLAGS) .
