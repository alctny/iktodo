.PHONY: start build

VERSION = $(shell git rev-parse --short HEAD)
APP     = iktodo

all: build

build:
	@go build -ldflags "-w -s -X main.Version=$(VERSION)" -o $(APP)

clean:
	@rm -f $(APP)
	@git ls-files --others --exclude-standard | xargs rm -rf
