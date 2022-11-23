SHELL = /bin/sh

APP_NAME ?= terraform-provider-ansiblevault
PACKAGES ?= ./...
GO_FILES ?= *.go */*/*.go
VERSION = $(shell git describe --tags)

GOBIN=bin
BINARY_PATH=$(GOBIN)/$(APP_NAME)

LIB_SOURCE = main.go

GO_ARCH ?= $(shell go env GOHOSTARCH)
GO_OS ?= $(shell go env GOHOSTOS)

TERRAFORM_PLUGIN_FOLDER ?= $(HOME)/.terraform.d/plugins

.DEFAULT_GOAL := app

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## name: Output name
.PHONY: name
name:
	@printf "%s" "$(APP_NAME)"

## version: Output sha1 of last commit
.PHONY: version
version:
	@printf "%s" "$(VERSION)"

##app: Build app with dependencies download
.PHONY: app
app: init dev

## dev: Build app
.PHONY: dev
dev: format style test build

## init: Download dependencies
.PHONY: init
init:
	go get github.com/kisielk/errcheck
	go get golang.org/x/lint/golint
	go install "golang.org/x/tools/cmd/goimports@latest"
	go mod tidy

## format: Format code
.PHONY: format
format:
	goimports -w $(GO_FILES)
	gofmt -s -w $(GO_FILES)

## style: Lint code
.PHONY: style
style:
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
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o "$(BINARY_PATH)_$(VERSION)" "$(LIB_SOURCE)"

# install: Install plugin into terraform plugin system
.PHONY: install
install:
	mkdir -p "$(TERRAFORM_PLUGIN_FOLDER)/$(GO_OS)_$(GO_ARCH)/"
	cp "$(BINARY_PATH)_$(VERSION)" "$(TERRAFORM_PLUGIN_FOLDER)/$(GO_OS)_$(GO_ARCH)/$(APP_NAME)_$(VERSION)"

# uninstall: Remove plugin from terraform plugin system
.PHONY: uninstall
uninstall:
	rm "$(TERRAFORM_PLUGIN_FOLDER)/$(GO_OS)_$(GO_ARCH)/$(APP_NAME)_$(VERSION)"

# clean: Delete binary
.PHONY: clean
clean:
	rm "$(BINARY_PATH)_$(VERSION)"

## github: Build and deploy on github
.PHONY: github
github:
	script/release

## config: Configure dev environment
.PHONY: config
config:
	script/install_hooks

.PHONY: all
all: build install
