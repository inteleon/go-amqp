all: test

test: deps
	@echo "Running tests..."
	@cd amqp && go fmt
	@cd amqp/queue && go fmt
	@cd monitor && go fmt
	@cd monitor/checker && go fmt
	@go test ./...

deps:
	@echo "Fetching dependencies..."
	@go get -u github.com/inteleon/go-logging/logging
	@go get -u github.com/inteleon/go-logging/helper
	@go get -u github.com/streadway/amqp
	@go get -u github.com/satori/go.uuid
