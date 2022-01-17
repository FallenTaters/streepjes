
wasm:
	GOARCH=wasm GOOS=js go build -o ./frontend/assets/app.wasm ./frontend/

run: wasm
	go run .

