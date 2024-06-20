.PHONY: run, build, test, lint

run: build
	./bin/app

run-qr: build-qr
	./bin/qr

run-auth: build-auth
	./bin/auth

build:
	go build -o ./bin/app cmd/server/main.go 

build-qr:
	go build -o ./bin/qr cmd/qr/main.go

build-auth:
	go build -o ./bin/auth cmd/auth/main.go
test:
	go test ./... -race

lint:
	golangci-lint run