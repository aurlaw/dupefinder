VERSION := $(shell git describe --tags --always --dirty)

build:
	go build -ldflags "-X github.com/aurlaw/dupefinder/cmd.Version=$(VERSION)" -o bin/dupefinder .

.PHONY: build