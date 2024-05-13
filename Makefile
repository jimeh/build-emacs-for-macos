PIP := $(shell command -v pip3 || command -v pip)
SOURCES := $(shell \
	find * \
	-not -path 'sources/*' \
	-not -path 'builds/*' \( \
		-name "*.go" -or \
		-name "go.mod" -or \
		-name "go.sum" -or \
		-name "Makefile" -or \
		-type f -path 'internal/*' -or \
		-type f -path 'cmd/*' -or \
		-type f -path 'pkg/*' \
	\) | grep -v '.DS_Store' \
)

#
# Environment
#

# Verbose output
ifdef VERBOSE
V = -v
endif

BINDIR := bin
TOOLDIR := $(BINDIR)/tools

# Global environment variables for all targets
SHELL ?= /bin/bash
SHELL := env \
	GO111MODULE=on \
	GOBIN=$(CURDIR)/$(BINDIR) \
	CGO_ENABLED=0 \
	PATH='$(CURDIR)/$(BINDIR):$(CURDIR)/$(TOOLDIR):$(PATH)' \
	$(SHELL)

#
# Defaults
#

# Default target
.DEFAULT_GOAL := build

#
# Bootstrap
#

bootstrap: bootstrap-brew bootstrap-ruby

bootstrap-ruby:
	bundle install

bootstrap-brew:
	brew bundle --verbose

bootstrap-pip:
	$(PIP) install -r requirements-ci.txt

#
# Tools
#

# external tool
define tool # 1: binary-name, 2: go-import-path
TOOLS += $(TOOLDIR)/$(1)

$(TOOLDIR)/$(1): Makefile
	GOBIN="$(CURDIR)/$(TOOLDIR)" go install "$(2)"
endef

$(eval $(call tool,gofumpt,mvdan.cc/gofumpt@latest))
$(eval $(call tool,golangci-lint,github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55))
$(eval $(call tool,gomod,github.com/Helcaraxan/gomod@latest))

.PHONY: tools
tools: $(TOOLS)

#
# Build
#

LDFLAGS := -w -s

VERSION ?= $(shell git describe --tags --exact 2>/dev/null)
COMMIT  ?= $(shell git rev-parse HEAD 2>/dev/null)
DATE    ?= $(shell date '+%FT%T%z')

ifeq ($(VERSION),)
	VERSION = 0.0.0-dev
endif

CMDDIR := cmd
BINS := $(shell test -d "$(CMDDIR)" && cd "$(CMDDIR)" && \
	find * -maxdepth 0 -type d -exec echo $(BINDIR)/{} \;)

.PHONY: build
build: $(BINS)

$(BINS): $(BINDIR)/%: $(SOURCES)
	mkdir -p "$(BINDIR)"
	cd "$(CMDDIR)/$*" && go build -a $(V) \
		-o "$(CURDIR)/$(BINDIR)/$*" \
		-ldflags "$(LDFLAGS) \
			-X main.version=$(VERSION) \
			-X main.commit=$(COMMIT) \
			-X main.date=$(DATE)"

#
# Development
#

TEST ?= $$(go list ./... | grep -v 'sources/' | grep -v 'builds/')

.PHONY: clean
clean:
	rm -rf $(BINARY) $(TOOLS)
	rm -f ./go.mod.tidy-check ./go.sum.tidy-check

.PHONY: test
test:
	CGO_ENABLED=1 go test $(V) -count=1 -race $(TESTARGS) $(TEST)

.PHONY: lint
lint: $(TOOLDIR)/golangci-lint
	golangci-lint $(V) run

.PHONY: format
format: $(TOOLDIR)/gofumpt
	gofumpt -w .

.PHONY: gen
gen:
	go generate $$(go list ./... | grep -v 'sources/' | grep -v 'builds/')

#
# Dependencies
#

.PHONY: deps
deps:
	$(info Downloading dependencies)
	go mod download

.PHONY: deps-update
deps-update:
	$(info Downloading dependencies)
	go get -u -t ./...

.PHONY: deps-analyze
deps-analyze: $(TOOLDIR)/gomod
	gomod analyze

.PHONY: tidy
tidy:
	go mod tidy $(V)

.PHONY: verify
verify:
	go mod verify

.SILENT: check-tidy
.PHONY: check-tidy
check-tidy:
	cp go.mod go.mod.tidy-check
	cp go.sum go.sum.tidy-check
	go mod tidy
	( \
		diff go.mod go.mod.tidy-check && \
		diff go.sum go.sum.tidy-check && \
		rm -f go.mod go.sum && \
		mv go.mod.tidy-check go.mod && \
		mv go.sum.tidy-check go.sum \
	) || ( \
		rm -f go.mod go.sum && \
		mv go.mod.tidy-check go.mod && \
		mv go.sum.tidy-check go.sum; \
		exit 1 \
	)
