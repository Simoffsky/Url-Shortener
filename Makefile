.PHONY: run, build, test, lint

run: build
	./bin/app

run-qr: build-qr
	./bin/qr

build:
	go build -o ./bin/app cmd/server/main.go 

build-qr:
	go build -o ./bin/qr cmd/qr/main.go

test:
	go test ./... -race

lint:
	golangci-lint run