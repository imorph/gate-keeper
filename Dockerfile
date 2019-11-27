FROM golang:1.13 as builder

RUN mkdir -p /gk/

WORKDIR /gk

COPY . .

RUN go mod download

RUN go test -v -race ./...

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
    -X github.com/imorph/gate-keeper/pkg/version.revision=${GIT_COMMIT}" \
    -a -o bin/gk cmd/gk/*

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
    -X github.com/imorph/gate-keeper/pkg/version.revision=${GIT_COMMIT}" \
    -a -o bin/gkcli cmd/gkcli/*

FROM alpine:3.10

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add \
    curl openssl netcat-openbsd

WORKDIR /home/app

COPY --from=builder /gk/bin/gk .
COPY --from=builder /gk/bin/gkcli /usr/local/bin/gkcli
RUN chown -R app:app ./

USER app

CMD ["./gk"]