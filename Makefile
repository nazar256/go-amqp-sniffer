doc: build
	bin/go-amqp-sniffer doc doc/

build: test
	go build -o bin/go-amqp-sniffer

vendor: tidy
	go mod vendor

tidy:
	go mod tidy

test: vendor
	golangci-lint run ./...
	go test -v  -race -timeout 5s ./...

lint:
	go vet ./...
	golangci-lint run ./... --fix

run:
	bin/go-amqp-sniffer  --exchange=test --parse-json | \
    		sudo docker exec -i go-amqp-sniffer-mongo \
    		mongoimport -d local -c test --batchSize 1 --type json

services:
	sudo docker-compose up -d
	sleep 15  # wait for rabbitmq to start

run-complete: services run