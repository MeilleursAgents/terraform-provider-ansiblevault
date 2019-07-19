SHELL = /bin/sh

APP_NAME ?= terraform-provider-ansiblevault
VERSION = v1.1.0
PACKAGES ?= ./...

GOBIN=bin
BINARY_PATH=$(GOBIN)/$(APP_NAME)

LIB_SOURCE = main.go

GO_ARCH=$(shell go env GOHOSTARCH)
GO_OS=$(shell go env GOHOSTOS)

.DEFAULT_GOAL := $(APP_NAME)

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## name: Output name
.PHONY: name
name:
	@echo -n $(APP_NAME)

## dist: Output binary path
.PHONY: dist
dist:
	@echo -n $(BINARY_PATH)

## version: Output sha1 of last commit
.PHONY: version
version:
	@echo -n $(VERSION)

## author: Output author's name of last commit
.PHONY: author
author:
	@python -c 'import sys; import urllib; sys.stdout.write(urllib.quote_plus(sys.argv[1]))' "$(shell git log --pretty=format:'%an' -n 1)"

## $(APP_NAME): Build app with dependencies download
.PHONY: $(APP_NAME)
$(APP_NAME): deps go

## go: Build app
.PHONY: go
go: format lint test build

## deps: Download dependencies
.PHONY: deps
deps:
	go get github.com/kisielk/errcheck
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports

## format: Format code
.PHONY: format
format:
	goimports -w *.go */*/*.go
	gofmt -s -w *.go */*/*.go

## lint: Lint code
.PHONY: lint
lint:
	golint $(PACKAGES)
	errcheck -ignoretests $(PACKAGES)
	go vet $(PACKAGES)

## test: Test with coverage
.PHONY: test
test:
	script/coverage

## build: Build binary
.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH)_$(VERSION) $(LIB_SOURCE)

.PHONY: install
install:
	mkdir -p $(HOME)/.terraform.d/plugins/$(GO_OS)_$(GO_ARCH)/
	cp $(BINARY_PATH)_$(VERSION) $(HOME)/.terraform.d/plugins/$(GO_OS)_$(GO_ARCH)/$(APP_NAME)_$(VERSION)

.PHONY: uninstall
uninstall:
	rm $(HOME)/.terraform.d/plugins/$(GO_OS)_$(GO_ARCH)/$(APP_NAME)_$(VERSION)

.PHONY: clean
clean:
	rm $(BINARY_PATH)_$(VERSION)

## github: build and deploy on github
.PHONY: github
github:
	go run cmd/release/*.go -access-token ${token}

.PHONY: all
all: build install
