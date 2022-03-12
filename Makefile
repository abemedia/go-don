all: lint test

lint:
	@golangci-lint run

todo:
	@golangci-lint run --no-config --disable-all --enable godox

test:
	@go test ./...

