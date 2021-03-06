## Overview
`rosetta-ethereum-2.0` provides a reference implementation of the Rosetta API for
Ethereum in Golang. If you haven't heard of the Rosetta API, you can find more
information [here](https://rosetta-api.org).

## Features
* Rosetta API implementation

## Usage
As specified in the [Rosetta API Principles](https://www.rosetta-api.org/docs/automated_deployment.html),
all Rosetta implementations must be deployable via Docker and support running via either an
[`online` or `offline` mode](https://www.rosetta-api.org/docs/node_deployment.html#multiple-modes).

**YOU MUST INSTALL DOCKER FOR THE FOLLOWING INSTRUCTIONS TO WORK. YOU CAN DOWNLOAD
DOCKER [HERE](https://www.docker.com/get-started).**

### Install
Running the following commands will create a Docker image called `rosetta-ethereum-2.0:latest`.

#### From Source
After cloning this repository, run:
```text
make build-local
```

### Run
Running the following commands will start a Docker container in
[detached mode](https://docs.docker.com/engine/reference/run/#detached--d) with
a data directory at `<working directory>/beacon-data` and the Rosetta API accessible
at port `8080`.

_It is possible to run `rosetta-ethereum-2.0` using a remote node by adding
`-e "BEACON_RPC=<node url>"` to any online command._

#### Mainnet:Online
```text
docker run -d --rm --ulimit "nofile=100000:100000" -v "$(pwd)/beacon-data:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -e "WEB3PROVIDER=<ETHEREUM URL>"  -p 8080:8080 -p 30303:30303 rosetta-ethereum-2.0:latest
```
_If you cloned the repository, you can run `make run-mainnet-online ethereumRPC=<NODE URL>`._

#### Mainnet:Online (Remote)
```text
docker run -d --rm --ulimit "nofile=100000:100000" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -e "WEB3PROVIDER=<ETHEREUM URL>"  -e "BEACON_RPC=<NODE URL>" -p 8080:8080 -p 30303:30303 rosetta-ethereum-2.0:latest
```
_If you cloned the repository, you can run `make run-mainnet-remote ethereumRPC=<NODE URL> beaconRPC=<NODE URL>`._

#### Mainnet:Offline
```text
docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -e "WEB3PROVIDER=<ETHEREUM URL>" -p 8081:8081 rosetta-ethereum-2.0:latest
```
_If you cloned the repository, you can run `make run-mainnet-offline ethereumRPC=<NODE URL>`._

#### Testnet:Online
```text
docker run -d --rm --ulimit "nofile=100000:100000" -v "$(pwd)/beacon-data:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -e "WEB3PROVIDER=<ETHEREUM URL>"  -p 8080:8080 -p 30303:30303 rosetta-ethereum-2.0:latest
```
_If you cloned the repository, you can run `make run-testnet-online ethereumRPC=<NODE URL>`._

#### Testnet:Online (Remote)
```text
docker run -d --rm --ulimit "nofile=100000:100000" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -e "WEB3PROVIDER=<ETHEREUM URL>"  -e "BEACON_RPC=<NODE URL>" -p 8080:8080 -p 30303:30303 rosetta-ethereum-2.0:latest
```
_If you cloned the repository, you can run `make run-testnet-remote ethereumRPC=<NODE URL> beaconRPC=<NODE URL>`._

#### Testnet:Offline
```text
docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -e "WEB3PROVIDER=<ETHEREUM URL>" -p 8081:8081 rosetta-ethereum-2.0:latest
```
_If you cloned the repository, you can run `make run-testnet-offline ethereumRPC=<NODE URL>`._

## System Requirements
`rosetta-ethereum-2.0` has been tested on an [AWS c5.2xlarge instance](https://aws.amazon.com/ec2/instance-types/c5).
This instance type has 8 vCPU and 16 GB of RAM. If you use a computer with less than 16 GB of RAM,
it is possible that `rosetta-ethereum-2.0` will exit with an OOM error.

### Recommended OS Settings
To increase the load `rosetta-ethereum-2.0` can handle, it is recommended to tune your OS
settings to allow for more connections. On a linux-based OS, you can run the following
commands ([source](http://www.tweaked.io/guide/kernel)):
```text
sysctl -w net.ipv4.tcp_tw_reuse=1
sysctl -w net.core.rmem_max=16777216
sysctl -w net.core.wmem_max=16777216
sysctl -w net.ipv4.tcp_max_syn_backlog=10000
sysctl -w net.core.somaxconn=10000
sysctl -p (when done)
```
_We have not tested `rosetta-ethereum-2.0` with `net.ipv4.tcp_tw_recycle` and do not recommend
enabling it._

You should also modify your open file settings to `100000`. This can be done on a linux-based OS
with the command: `ulimit -n 100000`.

## Testing with rosetta-cli
To validate `rosetta-ethereum-2.0`, [install `rosetta-cli`](https://github.com/coinbase/rosetta-cli#install)
and run one of the following commands:
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/testnet/config.json`
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/mainnet/config.json`

## Future Work
* [Rosetta API `/block/*`](https://www.rosetta-api.org/docs/BlockApi.html) add handling for missed blocks
* Add ERC-20 Rosetta Module to enable reading ERC-20 token transfers and transaction construction
* [Rosetta API `/mempool/*`](https://www.rosetta-api.org/docs/MempoolApi.html) implementation
* [Rosetta API `/account/*`](https://www.rosetta-api.org/docs/AccountApi.html) implementation
* [Rosetta API `/construction/*`](https://www.rosetta-api.org/docs/ConstructionApi.html) implementation
* Add CI test using `rosetta-cli` to run on each PR

## Development
* `make deps` to install dependencies
* `make test` to run tests
* `make build-local` to build a Docker image from the local context