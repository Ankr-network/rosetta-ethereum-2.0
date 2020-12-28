# Compile golang 
FROM ubuntu:18.04 as golang-builder

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app
WORKDIR /app

RUN apt-get update && apt-get install -y curl wget make gcc g++ git
ENV GOLANG_VERSION 1.15.5
ENV GOLANG_DOWNLOAD_SHA256 9a58494e8da722c3aef248c9227b0e9c528c7318309827780f16220998180a0d
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Get beacon
FROM golang-builder as beacon-builder

# VERSION: beacon v1.0.5
RUN wget --output-document beacon-chain https://github.com/prysmaticlabs/prysm/releases/download/v1.0.5/beacon-chain-v1.0.5-linux-amd64 \
  && mv beacon-chain /app/beacon-chain

# Compile rosetta-ethereum
FROM golang-builder as rosetta-builder

# Use native remote build context to build in any directory
COPY . src 
RUN cd src \
  && go build

RUN mv src/rosetta-ethereum-2.0 /app/rosetta-ethereum-2.0 \
  && mkdir /app/ethereum \
  && mv src/ethereum/prysm-config.yaml /app/ethereum/prysm-config.yaml \
  && rm -rf src 

## Build Final Image
FROM ubuntu:18.04

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app \
  && mkdir -p /data \
  && chown -R nobody:nogroup /data

WORKDIR /app

# Copy binary from beacon-builder
COPY --from=beacon-builder /app/beacon-chain /app/beacon-chain

# Copy binary from rosetta-builder
COPY --from=rosetta-builder /app/ethereum /app/ethereum
COPY --from=rosetta-builder /app/rosetta-ethereum-2.0 /app/rosetta-ethereum-2.0

# Set permissions for everything added to /app
RUN chmod -R 755 /app/*

CMD ["/app/rosetta-ethereum-2.0", "run"]
