
# BEGIN: lint-install ../infrastructure/
# http://github.com/tinkerbell/lint-install

GOLINT_VERSION ?= v1.42.0


LINT_OS := $(shell uname)
LINT_ARCH := $(shell uname -m)

# shellcheck and hadolint lack arm64 native binaries: rely on x86-64 emulation
ifeq ($(LINT_OS),Darwin)
	ifeq ($(LINT_ARCH),arm64)
		LINT_ARCH=x86_64
	endif
endif


GOLINT_CONFIG:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))/.golangci.yml

lint: out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH)
	find . -name go.mod | xargs -n1 dirname | xargs -n1 -I{} sh -c "cd {} && golangci-lint run -c $(GOLINT_CONFIG)"

out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH):
	mkdir -p out/linters
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b out/linters $(GOLINT_VERSION)
	mv out/linters/golangci-lint out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH)

# END: lint-install ../infrastructure/
