PROJECT_NAME := wotbot
PROJECT := github.com/L11R/wotbot
VERSION := $(shell cat version)
COMMIT := $(shell git rev-parse --short HEAD)
PKG_LIST := $(shell go list ./... | grep -v /vendor/)

GOLANGCI_LINT_VERSION = v1.21.0

LDFLAGS = "-s -w -X $(PROJECT)/internal/version.Version=$(VERSION) -X $(PROJECT)/internal/version.Commit=$(COMMIT)"

build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o ./bin/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)

test:
	@go test -v -cover -gcflags=-l --race $(PKG_LIST)

lint:
	@golangci-lint run -v
