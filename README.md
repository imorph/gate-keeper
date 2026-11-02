# gate-keeper

gate-keeper is simple anti-bruteforce service with gRPC API

```

                                                                     
           +----------+                            +-------------+                     +---------+
           |          |                            |             |   reset ip/login    |         |           
           |  Some    | ip+login+hash(pass) req    |             +<--------------------+         |
 Auth Req  |  Auth    +--------------------------->+             |    CIDR black list  |         |
+--------->+  Service |                            | Gate-Keeper +<--------------------+ gkcli   |
           |          +<---------------------------+             |    CIDR white list  |         |
<----------+          |       ok/nok resp          |             +<--------------------+         |
   ok/nok  |          |                            |             |                     |         |
           +----------+                            +-------------+                     +---------+



```

## Build

### Binaries

Requires Go 1.27.3 or newer. CI reads the Go version from `go.mod`; keep the
Docker builder version aligned when updating it.

```shell
make build
```

produces `gk` and `gkcli` in `./bin` directory

`gk` is service
`gkcli` is cli command able to:

* check if ip/login/pass banned or not
* add IP CIDR to white/black list
* reset tries counters for particular IP/Logins
* exec simple benchmark against service

### Code check

Check formatting, run `go vet`, golangci-lint (including errcheck), and govulncheck:

```shell
make check-all
```

The pinned check tools are installed into `bin/tools`. Run `make fmt` to apply
formatting, or `make vulncheck` to check for reachable known vulnerabilities.

### Unit tests

Tests include the race detector and local gRPC integration tests; no external
services are required.

```shell
make test
```

### Container

builds docker container:

```shell
make build-container
```

## Run

to build and run latest version in docker:

```shell
docker compose up -d
```

to run binaries natively:

```shell
make run
```

## Bench

to do unit benchmark `core`:

```shell
make bench
```

to do e2e benchmark `gkcli simple-bench` inside docker container:

```shell
make docker-bench
```

## TODO

* make white/black lists persitent somewhere (DB? file?)
* add option to move all state to external DB (Redis)
* add multicore/multi-instance benchmarks
* send notifications to Kafka topics
* add prometheus-style metrics
* move containers to Github registry
