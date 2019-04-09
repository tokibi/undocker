VERSION ?= $(shell git describe --tag --abbrev=0)

build: export GO111MODULE=on
build:
	go build --ldflags "-s -w -X main.version=$(VERSION)" ./cmd/undocker
