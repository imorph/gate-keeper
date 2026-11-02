# Makefile for releasing gate-keeper
#
# The release version is controlled from pkg/version

TAG?=latest
NAME:=gate-keeper
CLI_NAME:=gate-keeper-cli
DOCKER_REPOSITORY:=ivanvg
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
GIT_COMMIT:=$(shell git describe --dirty --always)
VERSION:=$(shell grep "version = " pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')

GOLANGCI_LINT_VERSION:=v2.15.0
GOVULNCHECK_VERSION:=v1.9.0
TOOLS_DIR:=$(CURDIR)/bin/tools
GOLANGCI_LINT:=$(TOOLS_DIR)/golangci-lint/$(GOLANGCI_LINT_VERSION)/golangci-lint
GOVULNCHECK:=$(TOOLS_DIR)/govulncheck/$(GOVULNCHECK_VERSION)/govulncheck
DOCKER_COMPOSE?=docker compose

.PHONY: run test bench docker-bench build release fmt fmt-check vet lint errcheck check-all golangci-lint install-golangci-lint vulncheck build-container push-container

run:
		GO111MODULE=on go run -ldflags "-s -w -X github.com/imorph/gate-keeper/pkg/version.revision=$(GIT_COMMIT)" cmd/gk/* --log-level=debug

test:
		GO111MODULE=on go test -v -race ./...

bench:
		GO111MODULE=on go test -bench=. ./pkg/core/... -benchmem
		GO111MODULE=on go test -bench=. ./pkg/server/... -benchmem

docker-bench:
		@set -eu; trap '$(DOCKER_COMPOSE) down' EXIT; \
		$(DOCKER_COMPOSE) up -d --build; \
		$(DOCKER_COMPOSE) exec -T gk gkcli simple-bench

build:
		GO111MODULE=on CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/imorph/gate-keeper/pkg/version.revision=$(GIT_COMMIT) -X github.com/imorph/gate-keeper/pkg/version.appname=$(NAME)" -a -o ./bin/gk ./cmd/gk/*
		GO111MODULE=on CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/imorph/gate-keeper/pkg/version.revision=$(GIT_COMMIT) -X github.com/imorph/gate-keeper/pkg/version.appname=$(CLI_NAME)" -a -o ./bin/gkcli ./cmd/gkcli/*

release: build-container push-container
		git tag $(VERSION)
		git push origin $(VERSION)
		goreleaser --rm-dist

fmt:
		GO111MODULE=on gofmt -l -w -s ./pkg
		GO111MODULE=on gofmt -l -w -s ./cmd

fmt-check:
		@test -z "$$(gofmt -l ./pkg ./cmd)" || { gofmt -l ./pkg ./cmd; exit 1; }

vet:
		GO111MODULE=on go vet ./pkg/...
		GO111MODULE=on go vet ./cmd/...

lint: golangci-lint

errcheck: $(GOLANGCI_LINT)
		$(GOLANGCI_LINT) run --default none --enable errcheck ./...

check-all: fmt-check vet lint vulncheck

golangci-lint: $(GOLANGCI_LINT)
		$(GOLANGCI_LINT) run ./...

install-golangci-lint: $(GOLANGCI_LINT)

$(GOLANGCI_LINT):
		GOBIN=$(dir $(GOLANGCI_LINT)) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

vulncheck: $(GOVULNCHECK)
		$(GOVULNCHECK) ./...

$(GOVULNCHECK):
		GOBIN=$(dir $(GOVULNCHECK)) go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

build-container:
		docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

push-container:
		docker push $(DOCKER_IMAGE_NAME):$(VERSION)
