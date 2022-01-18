all: build

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_TIME := $(shell git show -s --format=%ci $(GIT_COMMIT) | tr -d '\n')
GIT_TAG := $(shell git describe --exact-match --tags $(GIT_COMMIT) 2> /dev/null)

PACKAGE := github.com/PotatoesFall/vecty-test/src/infrastructure/router

LDFLAGS := '-X "$(PACKAGE).buildCommit=$(GIT_COMMIT)" \
		-X "$(PACKAGE).buildTime=$(GIT_TIME)" \
		-X "$(PACKAGE).buildVersion=$(GIT_TAG)"'

wasm:
	@GOARCH=wasm GOOS=js go build -o ./src/infrastructure/static/files/app.wasm ./frontend/

run: wasm
	@go run -ldflags "-X $(PACKAGE).buildVersion=development" .

build: wasm
	@go build -o ./bin/vecty-test -ldflags $(LDFLAGS) .
