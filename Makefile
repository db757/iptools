BINARY_NAME := ipt
DIST_DIR := ./dist
GOLIST := $(shell go list ./...)

build: tidy clean fmt vet test nix-update
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
	golint ${GOLIST}
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

nix-update:
	gomod2nix generate
.PHONY: nix-update

nix-build: nix-update
	nix build
.PHONY: nix-build

nix-install:
	nix profile install
.PHONY: nix-install

nix-shell:
	nix develop
.PHONY: nix-shell

targets: clean
	echo "Compiling targets"
	GOOS=windows GOARCH=386 go build -o ${DIST_DIR}/${BINARY_NAME}-windows-386.exe ./
	GOOS=windows GOARCH=amd64 go build -o ${DIST_DIR}/${BINARY_NAME}-windows-amd64.exe ./
	GOOS=linux GOARCH=arm64 go build -o ${DIST_DIR}/${BINARY_NAME}-linux-arm64 ./
	GOOS=linux GOARCH=amd64 go build -o ${DIST_DIR}/${BINARY_NAME}-linux-amd64 ./
	GOOS=darwin GOARCH=amd64 go build -o ${DIST_DIR}/${BINARY_NAME}-darwin-amd64 ./
	GOOS=darwin GOARCH=arm64 go build -o ${DIST_DIR}/${BINARY_NAME}-darwin-arm64 ./
.PHONY: targets

run: vet
	go run ./
.PHONY: run
