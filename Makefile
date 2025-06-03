BINARY_NAME := ipt
DIST_DIR := ./dist

build: tidy clean fmt vet test
	mkdir ${DIST_DIR}/
	go build -o ${DIST_DIR}/${BINARY_NAME} ./
.PHONY: build

tidy:
	go mod tidy
	go mod vendor
.PHONY: tidy

fmt:
	go fmt ./...
.PHONY: fmt

lint: fmt
	golint ./...
.PHONY: lint

test:
	go test ./...
.PHONY: test

vet: fmt
	go vet ./...
.PHONY: vet

govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck --show verbose ./...
.PHONY: govulncheck

upgrade:
	go get -u ./...
	go mod tidy
.PHONY: upgrade

clean:
	rm -Rf ${DIST_DIR}
	go clean
.PHONY: clean

targets: clean
	echo "Compiling targets"
	GOOS=linux GOARCH=amd64 go build -o ${DIST_DIR}/${BINARY_NAME}-linux-amd64 ./
	GOOS=darwin GOARCH=arm64 go build -o ${DIST_DIR}/${BINARY_NAME}-darwin-arm64 ./
.PHONY: targets

run: vet
	go run ./
.PHONY: run

