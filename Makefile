test:
	go test ./...

cover: test
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

public_html/assets/wasm_exec.js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" public_html/assets/

build: public_html/assets/wasm_exec.js
	cd lib && GOOS=js GOARCH=wasm go build -o ../public_html/assets/downsample.wasm
