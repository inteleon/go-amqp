all: test

test: deps
	@echo "Running tests..."
	@cd amqp && go fmt
	@go test ./...

deps:
	@echo "Fetching dependencies..."
	@go get -u github.com/inteleon/go-logging/logging
	@go get -u github.com/streadway/amqp
