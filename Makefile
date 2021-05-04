release: test doc build-windows-amd64 build-linux-amd64 build-linux-arm build-darwin-amd64 build-darwin-arm64

doc: build
	bin/go-amqp-sniffer doc doc/

build: test
	go build -o bin/go-amqp-sniffer

build-windows-amd64:
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
	zip bin/go-amqp-sniffer-windows-amd64.zip go-amqp-sniffer.exe
	go clean

build-linux-amd64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	zip bin/go-amqp-sniffer-linux-amd64.zip go-amqp-sniffer
	go clean

build-linux-arm:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build
	zip bin/go-amqp-sniffer-linux-arm.zip go-amqp-sniffer
	go clean

build-darwin-arm64:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build
	zip bin/go-amqp-sniffer-darwin-arm64.zip go-amqp-sniffer
	go clean

build-darwin-amd64:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
	zip bin/go-amqp-sniffer-darwin-amd64.zip go-amqp-sniffer
	go clean

vendor: tidy
	go mod vendor

tidy:
	go mod tidy

test: vendor
	golangci-lint run ./...
	go test -v  -race -timeout 5s ./...

lint:
	golangci-lint run ./... --fix

run:
	bin/go-amqp-sniffer  --exchange=test --parse-json | \
    		sudo docker exec -i go-amqp-sniffer-mongo \
    		mongoimport -d local -c test --batchSize 1 --type json

services:
	sudo docker-compose up -d
	sleep 15  # wait for rabbitmq to start

run-complete: services run