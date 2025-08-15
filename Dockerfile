# Compile application w/required OS library dependencies
FROM golang:1.23 AS builder
RUN apt-get update && \
    apt-get install -y gcc-aarch64-linux-gnu musl-dev g++-x86-64-linux-gnu libc6-dev-amd64-cross && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 go build -o protun-app main.go

# Run executable w/requires OS utilities (openvpn, etc)
FROM ubuntu:22.04
RUN apt-get update && \
    apt-get install -y openvpn dnsutils && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/protun-app .

# need to mount OS source to /vpn/configs when running
COPY /vpn/configs /app/ovpn/configs
RUN chmod +x /app/protun-app && chmod -R 755 /app/ovpn
CMD ["./protun-app"]


