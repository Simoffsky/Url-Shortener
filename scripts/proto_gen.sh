#!/bin/bash


cd "$(dirname "$0")/.."


mkdir -p ./pkg/proto/qr
mkdir -p ./pkg/proto/auth
mkdir -p ./pkg/proto/stats
protoc  --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./pkg/proto/qr/qr.proto

protoc  --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./pkg/proto/auth/auth.proto


protoc  --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./pkg/proto/stats/stats.proto