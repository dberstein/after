build_bin = after
version := $(shell git describe --tags --dirty)-$(shell git rev-parse HEAD)
targets := $(shell go tool dist list -json | jq -r '.[]|select(.FirstClass==true)|"\(.GOOS)/\(.GOARCH)"')
runtime := $(shell go env GOHOSTOS)/$(shell go env GOHOSTARCH)

ifeq ($(shell uname -s),Darwin)
install_dir := /usr/local/bin
else
install_dir := /usr/bin
endif

.PHONY: targets
targets:
	@echo "make all"
	@echo $(targets)|tr ' ' '\n'|sed 's/^/ make /g'|sed -r "s# make $(runtime)#(make $(runtime) > make install > make uninstall)#g"

.PHONY: all
all: $(targets)

.PHONY: $(targets)
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
$(targets): %: test
	@GOOS=$(os) GOARCH=$(arch) \
	go build -x \
	-ldflags="-X main.Version=$(version) -extldflags=-static" \
	-o 'build/$(os)/$(arch)/$(version)/$(build_bin)$(shell [ "$(os)" != "windows" ] || echo .exe)' \
	.

.PHONY: test
test:
	@go test -v -cover ./...

.PHONY: install
install: $(runtime)
	@install -v -b -m 755 "build/$(runtime)/$(version)/$(build_bin)" $(install_dir)

.PHONY: uninstall
uninstall:
	@rm -v $(install_dir)/$(build_bin)

.PHONY: clear
clean:
	@rm -rvf ./build
