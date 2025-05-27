build:
	@go build -o bin/ddfss

run: build
	@./bin/ddfss

test:
	@go test ./... -v
