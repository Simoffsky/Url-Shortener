.PHONY: run, build, test, lint

run: build
	./bin/bot

build:
	go build -o ./bin/bot cmd/bot/main.go 

test:
	go test ./... -race

lint:
	golangci-lint run