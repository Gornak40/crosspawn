.PHONY: build run clean test

SERVICE_NAME=crosspawn

BINARY_PATH=./bin/${SERVICE_NAME}

build:
	go build -o $(BINARY_PATH) ./cmd/$(SERVICE_NAME)

run: build
	@$(BINARY_PATH)

clean:
	rm -f $(BINARY_PATH)

test:
	go test -v ./...

lint:
	golangci-lint run ./...
