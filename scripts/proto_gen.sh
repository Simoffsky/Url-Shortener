#!/bin/bash


cd "$(dirname "$0")/.."


mkdir -p ./pkg/proto/qr

protoc  --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./pkg/proto/qr/qr.proto