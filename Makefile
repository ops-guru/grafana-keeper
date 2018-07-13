# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=grafana-keeper

# These will be provided to the target
REPO?=opsguru.io/$(BINARY_NAME)
VERSION := 0.0.1
BUILD?=$(shell git rev-parse --short HEAD)
TAG=$(VERSION)-$(BUILD)

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

all: build

clean:
	$(GOCLEAN)

build: clean
	$(GOBUILD) -o $(BINARY_NAME) -v $(LDFLAGS)

image: build
	docker build -t $(REPO):$(TAG) .

fmt:
	gofmt -d -e -l -s *.go
	gofmt -d -e -l -s keeper/*.go

golint:
	golint #
	golint keeper