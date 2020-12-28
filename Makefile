.PHONY: deps build run lint run-mainnet-online run-mainnet-offline run-testnet-online \
	run-testnet-offline build-local 

GO_PACKAGES=./services/... 
GO_FOLDERS=$(shell echo ${GO_PACKAGES} | sed -e "s/\.\///g" | sed -e "s/\/\.\.\.//g")
TEST_SCRIPT=go test ${GO_PACKAGES}
PWD=$(shell pwd)
NOFILE=100000

deps:
	go get ./...

test:
	${TEST_SCRIPT}

build:
	docker build -t rosetta-ethereum-2.0:latest https://github.com/Ankr-network/rosetta-ethereum-2.0.git

build-local:
	docker build -t rosetta-ethereum-2.0:latest .

build-release:
	# make sure to always set version with vX.X.X
	docker build -t rosetta-ethereum-2.0:$(version) .;
	docker save rosetta-ethereum-2.0:$(version) | gzip > rosetta-ethereum-2.0-$(version).tar.gz;

run-mainnet-online:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/beacon-data:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -e "WEB3PROVIDER=$(ethereumRPC)" -p 8080:8080 -p 30303:30303 rosetta-ethereum-2.0:latest

run-mainnet-offline:
	docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -e "WEB3PROVIDER=$(ethereumRPC)" -p 8081:8081 rosetta-ethereum-2.0:latest

run-testnet-online:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/beacon-data:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -e "WEB3PROVIDER=$(ethereumRPC)" -p 8080:8080 -p 30303:30303 rosetta-ethereum-2.0:latest

run-testnet-offline:
	docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -e "WEB3PROVIDER=$(ethereumRPC)" -p 8081:8081 rosetta-ethereum-2.0:latest

run-mainnet-remote:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -e "WEB3PROVIDER=$(ethereumRPC)" -e "BEACON_RPC=$(beaconRPC)" -p 8080:8080 -p 30303:30303 rosetta-ethereum-2.0:latest

run-testnet-remote:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -e "WEB3PROVIDER=$(ethereumRPC)" -e "BEACON_RPC=$(beaconRPC)" -p 8080:8080 -p 30303:30303 rosetta-ethereum-2.0:latest
