.PHONY: build test clean

build:
	go build -v ./...

test:
	go test -v ./...

clean:
	go clean
	rm -f goflipdot

run-example:
	go run cmd/example/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint ./...

all: fmt vet lint build test
