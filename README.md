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

VETing/linting/errchecking:

```shell
make check-all
```

### Unit tests

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
docker-compose up -d
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
