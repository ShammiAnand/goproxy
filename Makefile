GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=goproxy
BINARY_UNIX=$(BINARY_NAME)_unix
CONFIG_PATH=./internal/config/config.yaml

MAIN_PACKAGE=./cmd/goproxy

UNIT_TEST_PATH=./tests/unit/...
INTEGRATION_TEST_PATH=./tests/integration/...
PERFORMANCE_TEST_PATH=./tests/performance/...
SIMULATION_TEST_PATH=./tests/simulation/...

.PHONY: all build clean test run deps help

all: test build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

test: unit-test integration-test

unit-test:
	$(GOTEST) -v $(UNIT_TEST_PATH)

integration-test:
	$(GOTEST) -v $(INTEGRATION_TEST_PATH)

performance-test:
	$(GOTEST) -v -bench=. $(PERFORMANCE_TEST_PATH)

simulation:
	$(GOTEST) -v $(SIMULATION_TEST_PATH)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)
	./$(BINARY_NAME) -config $(CONFIG_PATH)

deps:
	$(GOGET) -v -d ./...

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PACKAGE)

docker-build:
	docker build -t $(BINARY_NAME):latest .

help:
	@echo "Make targets:"
	@echo "  build           - Build the GoProxy binary"
	@echo "  clean           - Remove binary files and clean up"
	@echo "  test            - Run all tests (unit, integration, performance)"
	@echo "  unit-test       - Run unit tests"
	@echo "  integration-test- Run integration tests"
	@echo "  performance-test- Run performance tests"
	@echo "  run             - Build and run the GoProxy"
	@echo "  deps            - Get dependencies"
	@echo "  build-linux     - Cross-compile for Linux"
	@echo "  docker-build    - Build Docker image"
	@echo "  help            - Display this help message"
