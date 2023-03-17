BUILD_BIN=after
BUILD_PREFIX=build/$(shell uname -s)-$(shell uname -m)
INSTALL_PREFIX=/usr/local

SOURCES=go.mod $(shell find . -type f -name '*.go')

#.PHONY: build
build: $(SOURCES)
	@go build -ldflags="-extldflags=-static" \
		-o $(BUILD_PREFIX)/$(BUILD_BIN) \
	cmd/after.go \
	&& strip $(BUILD_PREFIX)/$(BUILD_BIN)

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
	@install -m 755 $(BUILD_PREFIX)/$(BUILD_BIN) $(INSTALL_PREFIX)/bin
