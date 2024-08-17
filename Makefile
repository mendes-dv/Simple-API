# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=myapp
SRC_DIRS=$(shell find . -type d -print)

# Define build environment
.PHONY: all dep build clean run

# Default target executed when no arguments are given to make
all: build

# Resolves dependencies
dep:
	$(GOGET) ./...

# Builds the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

# Runs the project
run:
	$(GORUN) main.go

# Cleans up generated files
clean:
	go clean
	rm -f $(BINARY_NAME)