GOFILES=$(shell find . -type f -name '*.go')

clitest:
	@-rm *.jpg
	go run .

test:
	go test downsample/pkg/...
	cd cmd/wasm && GOOS=js GOARCH=wasm go test ./...

cover: test
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

public_html/assets/wasm_exec.js: Makefile
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" public_html/assets/

public_html/assets/downsample.wasm: Makefile $(GOFILES)
	cd cmd/wasm && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ../../public_html/assets/downsample.wasm

build: public_html/assets/wasm_exec.js public_html/assets/downsample.wasm Makefile
	@echo "Done"

run: build
	go run downsample/cmd/server


bench:
	go test -bench=. -run=^Nope_ -benchmem

profile:
	go test -bench=. -run=^Nope_ -memprofile=mem.prof -cpuprofile=cpu.prof -benchtime=10s

profile-memory: profile
	go tool pprof mem.prof

profile-cpu: profile
	go tool pprof cpu.prof
