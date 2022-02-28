all: lint test

lint:
	@go vet ./...

todo:
	@golangci-lint run --no-config --disable-all --enable godox

test:
	@go test ./...

