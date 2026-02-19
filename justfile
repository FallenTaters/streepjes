set dotenv-filename := ".env.dev"
set dotenv-load

git_commit := `git rev-parse --short HEAD`
git_time := `git show -s --format=%ci HEAD | tr -d '\n'`
git_tag := `git describe --exact-match --tags HEAD 2>/dev/null || true`

package := "github.com/FallenTaters/streepjes/backend/infrastructure/router"

ldflags := '-X "' + package + '.buildCommit=' + git_commit + '" -X "' + package + '.buildTime=' + git_time + '" -X "' + package + '.buildVersion=' + git_tag + '"'

default: build

generate:
    go generate ./...

vugugen:
    bash ./frontend/generate.bash

wasm:
    GOARCH=wasm GOOS=js go build -o ./static/files/app.wasm ./frontend/

run: vugugen wasm
    go run -ldflags "-X {{package}}.buildVersion=development" .

run-backend:
    go run -ldflags "-X {{package}}.buildVersion=development" .

test:
    go test ./backend/... -cover

lint:
    golangci-lint run ./backend/...

build: generate vugugen wasm
    CGO_ENABLED=1 go build -o ./bin/streepjes -ldflags '{{ldflags}}' .

build-arm: generate vugugen wasm
    GOOS=linux GOARCH=arm CGO_ENABLED=1 go build -o ./bin/streepjes-linux-arm -ldflags '{{ldflags}}' .

container:
    podman build \
        --build-arg GIT_COMMIT={{git_commit}} \
        --build-arg GIT_TIME="{{git_time}}" \
        --build-arg GIT_TAG="{{git_tag}}" \
        -t streepjes .
