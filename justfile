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

run:
    find . -name '*.go' -o -name '*.html' -o -name '*.js' | entr -cr go run -ldflags "-X {{package}}.buildVersion=development" .

run-once:
    go run -ldflags "-X {{package}}.buildVersion=development" .

test:
    go test ./backend/... -cover

lint:
    golangci-lint run ./backend/...

build: generate
    CGO_ENABLED=0 go build -o ./bin/streepjes -ldflags '{{ldflags}}' .

container:
    podman build \
        --build-arg GIT_COMMIT={{git_commit}} \
        --build-arg GIT_TIME="{{git_time}}" \
        --build-arg GIT_TAG="{{git_tag}}" \
        -t streepjes .
