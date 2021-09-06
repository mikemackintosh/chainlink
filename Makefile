all: test build

test:
	go test -v ./...

build:
	go build -o bin/chainlink cmd/chainlink/*.go

start:
	go run cmd/chainlink/*.go
