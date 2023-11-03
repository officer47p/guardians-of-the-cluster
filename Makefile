build:
	@go build -o ./bin/rate-limiter .

start: build
	@./bin/rate-limiter

test:
	@go test -v ./...