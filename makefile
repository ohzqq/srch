build:
	GOARCH=wasm GOOS=js go build -o wasm/www/assets/srch.wasm wasm/main.go
