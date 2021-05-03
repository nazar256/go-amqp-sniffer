# go-amqp-sniffer
AMQP exchange sniffer written in golang. Subscribes for events and dumps them in requested format (csv, json) to stdout.

# Use cases
* debugging or investigating messages in RabbitMQ-based event system (mongoimport, ELK, etc)
* event logging for AMQP

# Installation

Requires installed go 1.16 or later
```bash
go install github.com/nazar256/go-amqp-sniffer@latest
```

# Usage
See [usage documentation](doc/go-amqp-sniffer.md).
