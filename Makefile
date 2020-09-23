ifeq ($(SHELL), cmd)
	VERSION := $(shell git describe --exact-match --tags 2>nil)
	HOME := $(HOMEPATH)
else ifeq ($(SHELL), sh.exe)
	VERSION := $(shell git describe --exact-match --tags 2>nil)
	HOME := $(HOMEPATH)
else
	VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
endif

PREFIX := /usr
BINDIR := $(PREFIX)/bin
ETCDIR := /etc
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git rev-parse --short HEAD)
GOFILES ?= $(shell git ls-files '*.go')
GOFMT ?= $(shell gofmt -l -s $(filter-out plugins/parsers/influx/machine.go, $(GOFILES)))
BUILDFLAGS ?=
INSTALL := install
INSTALL_EXEC := $(INSTALL) -D --mode 755
INSTALL_DATA := $(INSTALL) -D --mode 0644

ifdef GOBIN
PATH := $(GOBIN):$(PATH)
else
PATH := $(subst :,/bin:,$(shell go env GOPATH))/bin:$(PATH)
endif

LDFLAGS := -X github.com/sunet/s3-mm-tool/pkg/meta.commit=$(COMMIT) -X github.com/sunet/s3-mm-tool/pkg/meta.branch=$(BRANCH)
ifdef VERSION
	LDFLAGS += -X github.com/sunet/s3-mm-tool/pkg/meta.version=$(VERSION)
endif

.PHONY: all
all:
	@$(MAKE) --no-print-directory s3-mm-tool #docs/mix.1

.PHONY: mix
mix:
	go build $(GO_BUILD_FLAGS) -ldflags "$(LDFLAGS)" ./cmd/tool

docs/%.1: docs/%.ronn.1
	ronn -r $< > $@

.PHONY: install
install: s3-mm-tool
	$(INSTALL_EXEC) s3-mm-tool $(DESTDIR)$(BINDIR)


.PHONY: test
test:
	go test $(GO_BUILD_FLAGS) -cover -short ./...

.PHONY: testcover
testcover:
	go test $(GO_BUILD_FLAGS) -cover ./...

.PHONY: test-all
test-all: fmtcheck vet
	go test ./...

.PHONY: clean
clean:
	rm -f s3-mm-tool
	rm -f s3-mm-tool.exe

.PHONY: docker
docker:
	docker build -t "s3-mm-tool:$(COMMIT)" .
	docker tag s3-mm-tool:$(COMMIT) docker.sunet.se/s3-mm-tool:$(COMMIT)
	docker push docker.sunet.se/s3-mm-tool:$(COMMIT)
