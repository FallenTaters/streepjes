all: build

wasm:
	@GOARCH=wasm GOOS=js go build -o ./src/infrastructure/static/files/app.wasm ./frontend/

run: wasm
	@go run .

build: wasm
	@go build -o ./bin/vecty-test .
