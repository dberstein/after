BUILD_BIN=after
BUILD_DIR=build/$(shell uname -s)-$(shell uname -m)
INSTALL_DIR=/usr/bin

SOURCES=go.mod $(shell find . -type f -name '*.go')

.PHONY: test
test:
	@go test -v ./cmd

.PHONY: build
build: test $(SOURCES)
	@go build -x -ldflags="-extldflags=-static" \
		-o $(BUILD_DIR)/$(BUILD_BIN) \
	cmd/after.go \
	&& strip $(BUILD_DIR)/$(BUILD_BIN)

.PHONY: build/docker
build/docker: build/image
	@docker run --rm -v $(PWD):/app \
		--entrypoint make \
			after:builder \
		-C /app build

build/image: Dockerfile
	@docker build -t after:builder .

.PHONY: install
install: build
	@install -v -m 755 $(BUILD_DIR)/$(BUILD_BIN) $(INSTALL_DIR)

.PHONY: uninstall
uninstall:
	@rm -v $(INSTALL_DIR)/$(BUILD_BIN)

.PHONY: clear
clear:
	@rm -rvf ./build
