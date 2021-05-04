go-amqp-sniffer
===============
AMQP exchange sniffer written in golang. Subscribes for events and dumps them in requested format (csv, json) to stdout.

Have you ever wondered what messages are passing by in your RabbitMQ?
If you have an event system bus build on top of AMQP broker - 
at some point you may lose track of what events are produced from or where the produced events are consumed.
This problem evolves when the event bus is shared between different teams.
In this case of investigating event such a sniffer may become handy, as you can attach to the bus and listen to the events
without any harm to other services (except maybe only additional load to the AMQP server).
Moreover, this tool can be used for logging the events directly to mongodb using 

Use cases
---------
* debugging messages in RabbitMQ-based event system (stream to mongoimport, ELK, etc)
* event logging for AMQP for case investigation

Installation
------------

### From sources (recommended)

Requires installed go 1.16 or later
```bash
go install github.com/nazar256/go-amqp-sniffer@latest
```

### From pre-compiled binaries

Check out the [release page](https://github.com/nazar256/go-amqp-sniffer/releases), there you may find a zipped binary for your platform.

Features
--------
* subscribes to AMQP exchanges using routing keys
* creates auto-removable (if not specified `--persistent-queue`) queue so when sniffer is off, queue will not be overfilled
* streams received messages to standard output as mongoimport compatible JSON or CSV
* input message bodies can be parsed as JSON and included in result as sub-object (can be saved as sub-documents in MongoDB)

Usage
-----
See [usage documentation](doc/go-amqp-sniffer.md).
