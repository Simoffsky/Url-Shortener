.PHONY: run, build, test, lint

run: build
	./bin/app

build:
	go build -o ./bin/app cmd/server/main.go 

test:
	go test ./... -race

lint:
	golangci-lint run