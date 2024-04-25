.PHONY: build run clean test lint

BINARIES=crosspawn poller
BIN_FOLDER=bin

build: $(BINARIES)

crosspawn:
	@go build -o ${BIN_FOLDER}/crosspawn ./cmd/crosspawn

poller:
	@go build -o ${BIN_FOLDER}/poller ./cmd/poller

run-crosspawn: crosspawn
	@${BIN_FOLDER}/crosspawn

run-poller: poller
	@${BIN_FOLDER}/poller

clean:
	@rm $(BIN_FOLDER)/*

test:
	@go test -v ./...

lint:
	@golangci-lint run ./...
