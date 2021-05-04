package broker

import (
	"fmt"

	"github.com/nazar256/go-amqp-sniffer/flow"
	"github.com/streadway/amqp"
)

type AMQPMsgChan <-chan amqp.Delivery

type Config struct {
	URL             string
	Prefetch        int
	RoutingKeys     []string
	QueueName       string
	ExchangeName    string
	PersistentQueue bool
}

type Listener struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	name string
}

func NewListener(cfg *Config) *Listener {
	if cfg.Prefetch < 1 {
		cfg.Prefetch = 1
	}

	conn, err := amqp.Dial(cfg.URL)
	flow.FailOnError(err, "Failed to connect to AMQP")

	ch, err := conn.Channel()
	flow.FailOnError(err, "Could not create amqp channel")

	err = ch.Qos(cfg.Prefetch, 0, false)
	flow.FailOnError(err, "Could not set prefetch")

	queue, err := ch.QueueDeclare(
		cfg.QueueName,
		cfg.PersistentQueue,  // durable
		!cfg.PersistentQueue, // delete when unused
		!cfg.PersistentQueue, // exclusive
		false,                // noWait
		nil,                  // arguments
	)
	flow.FailOnError(err, "Error declaring the Queue")

	for _, routingKey := range cfg.RoutingKeys {
		err = ch.QueueBind(
			queue.Name,       // name of the queue
			routingKey,       // bindingKey
			cfg.ExchangeName, // sourceExchange
			false,            // noWait
			nil,              // arguments
		)

		flow.FailOnError(err, fmt.Sprintf(
			"Error binding exchange %s to the queue %s with routing key %s",
			cfg.ExchangeName,
			queue.Name,
			routingKey,
		))
	}

	return &Listener{conn, ch, queue.Name}
}

func (l Listener) Listen() AMQPMsgChan {
	consumerCh, err := l.ch.Consume(
		l.name, // queue
		l.name, // consumer
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	flow.FailOnError(err, "Could not start consumer")

	return consumerCh
}

func (l Listener) Cancel() error {
	err := l.ch.Cancel(l.name, false)
	if err != nil {
		return err
	}

	return l.conn.Close()
}
