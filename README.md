go-amqp-sniffer
===============

[![tests](https://github.com/nazar256/go-amqp-sniffer/actions/workflows/tests.yml/badge.svg)](https://github.com/nazar256/go-amqp-sniffer/actions/workflows/tests.yml)
[![linters](https://github.com/nazar256/go-amqp-sniffer/actions/workflows/linters.yml/badge.svg)](https://github.com/nazar256/go-amqp-sniffer/actions/workflows/linters.yml)


AMQP exchange sniffer written in golang. Subscribes for events and dumps them in requested format (csv, json) to stdout.

Have you ever wondered what messages are passing by in your RabbitMQ?
If you have an event system bus build on top of AMQP broker - 
at some point you may lose track of what events are produced from or where the produced events are consumed.
This problem evolves when the event bus is shared between different teams.
In this case for investigating events such a sniffer may become handy, as you can attach to the bus and listen to the events
without any harm to other services (except only additional load to AMQP server).
Moreover, this tool can be used for logging the events directly to mongodb using shell pipe.

Use cases
---------
* debugging messages in RabbitMQ-based event system
* event logging for AMQP for case investigation (stream to mongoimport, ELK, etc)

Installation
------------

### From sources (recommended)

Using go 1.16 or later
```bash
go install github.com/nazar256/go-amqp-sniffer@latest
```

Using go prior to 1.16 (not tested)
```bash
go install github.com/nazar256/go-amqp-sniffer
```

### From pre-compiled binaries

Check out the [releases page](https://github.com/nazar256/go-amqp-sniffer/releases), 
there you may find a zipped binary for your platform.

Features
--------
* subscribes to AMQP exchanges using routing keys
* creates auto-removable (if not specified `--persistent-queue`) queue, so queue cannot be overfilled when a sniffer is off
* streams received messages to standard output as mongoimport compatible JSON or CSV
* input message bodies can be parsed as JSON and included in result as sub-object (can be saved as sub-documents in MongoDB)

Usage
-----
See [usage documentation](doc/go-amqp-sniffer.md).
