build:
	GOARCH=wasm GOOS=js go build -o testdata/www/assets/srch.wasm wasm/main.go
