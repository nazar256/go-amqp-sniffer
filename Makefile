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
	golangci-lint run ./... --fix

run: build
	sudo docker-compose up -d
	bin/go-amqp-sniffer  --exchange=test --prefetch=5 --routing-key="#" --parse-json | \
		sudo docker exec -i go-amqp-sniffer-mongo \
		mongoimport -d local -c test --type json
	sudo docker-compose stop