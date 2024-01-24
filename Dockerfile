# First stage: Rust build using the official Rust image
FROM rust:latest as rust-builder

# Set your working directory for Rust components
WORKDIR /usr/src/app

# Copy the Rust source code
COPY crypto/rust ./crypto/rust
COPY zk/rust ./zk/rust

# Build Rust components
RUN cd crypto/rust && cargo build --release && \
    cd ../../zk/rust && cargo build --release

# Second stage: Go build using the official Go image
FROM golang:1.21 as go-builder

# Install OpenSSL development libraries
RUN apt-get update && apt-get install -y \
    build-essential \
    libssl-dev \
    pkg-config

# Set your working directory for Go components
WORKDIR /go/src/app

# Copy the Go source code
COPY . .

# Copy the Rust binaries from the previous stage
COPY --from=rust-builder /usr/src/app/crypto/rust/target/release /usr/src/app/crypto/rust/target/release
COPY --from=rust-builder /usr/src/app/zk/rust/target/release /usr/src/app/zk/rust/target/release

# Create the directory if it doesn't exist
RUN mkdir -p /go/bin

# Build Go components
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN go build -o /go/bin/ilxd
RUN cd cli && go build -o /go/bin/ilxcli

# Final stage: setup the runtime environment
FROM ubuntu:22.04

# Copy the built binaries from previous stages
COPY --from=go-builder /go/bin/ilxd /usr/local/bin/ilxd
COPY --from=go-builder /go/bin/ilxcli /usr/local/bin/ilxcli

# Copy the entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 9002
EXPOSE 5001
ENV PATH="/usr/local/bin:${PATH}"

# Set the entrypoint to a shell
ENTRYPOINT ["/bin/sh", "-c"]

# Default command
CMD ["ilxd"]
