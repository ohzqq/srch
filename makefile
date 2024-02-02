build:
	GOARCH=wasm GOOS=js go build -o testdata/srch.wasm wasm/main.go
serve: build
	go run internal/srv/main.go
