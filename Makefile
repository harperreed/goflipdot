.PHONY: build test clean

build:
	go build -v ./...

build-arm:
	GOARCH=arm GOARM=7 GOOS=linux go build -v ./...

cli:
	GOARCH=arm GOARM=7 GOOS=linux go build -o flipdot-cli cmd/flipdot-cli/main.go

example:
	GOARCH=arm GOARM=7 GOOS=linux go build -o flipdot-example cmd/example/main.go



test:
	go test -v ./...

clean:
	go clean
	rm -f goflipdot

run-example:
	go run cmd/example/main.go cmd/example/patterns.go $(ARGS)

run-cli:
	go run cmd/flipdot-cli/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint ./...

all: fmt vet lint build test
