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

run:
		GO111MODULE=on go run -ldflags "-s -w -X github.com/imorph/gate-keeper/pkg/version.revision=$(GIT_COMMIT)" cmd/gk/* --log-level=debug

test:
		GO111MODULE=on go test -v -race ./...

bench:
		GO111MODULE=on go test -bench=. ./pkg/core/... -benchmem
		GO111MODULE=on go test -bench=. ./pkg/server/... -benchmem

docker-bench:
		docker-compose up -d
		docker exec gate-keeper_gk_1 gkcli simple-bench
		docker-compose down

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

vet:
		GO111MODULE=on go vet ./pkg/...
		GO111MODULE=on go vet ./cmd/...

lint: install-golint
		golint pkg/...
		golint cmd/...

install-golint:
		which golint || GO111MODULE=off go get -u golang.org/x/lint/golint

errcheck: install-errcheck
		errcheck ./pkg/...
		errcheck ./cmd/...

install-errcheck:
		which errcheck || GO111MODULE=off go get -u github.com/kisielk/errcheck

check-all: fmt vet lint errcheck golangci-lint

golangci-lint: install-golangci-lint
		golangci-lint --version
		golangci-lint run -D errcheck -D structcheck

install-golangci-lint:
		which golangci-lint || GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

build-container:
		docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

push-container:
		docker push $(DOCKER_IMAGE_NAME):$(VERSION)
