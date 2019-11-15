# Makefile for releasing gate-keeper
#
# The release version is controlled from pkg/version

TAG?=latest
NAME:=gate-keeper
CLI_NAME:=gate-keeper-cli
DOCKER_REPOSITORY:=imorph
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
GIT_COMMIT:=$(shell git describe --dirty --always)
VERSION:=$(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')

run:
		GO111MODULE=on go run -ldflags "-s -w -X github.com/imorph/gate-keeper/pkg/version.REVISION=$(GIT_COMMIT)" cmd/gk/* --log-level=debug

test:
		GO111MODULE=on go test -v -race ./...

build:
		GO111MODULE=on CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/imorph/gate-keeper/pkg/version.REVISION=$(GIT_COMMIT) -X github.com/imorph/gate-keeper/pkg/version.APPNAME=$(NAME)" -a -o ./bin/gk ./cmd/gk/*
		GO111MODULE=on CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/imorph/gate-keeper/pkg/version.REVISION=$(GIT_COMMIT) -X github.com/imorph/gate-keeper/pkg/version.APPNAME=$(CLI_NAME)" -a -o ./bin/gkcli ./cmd/gkcli/*

release:
		git tag $(VERSION)
		git push origin $(VERSION)
		goreleaser