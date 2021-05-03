package sniffer_test

import (
	"time"

	"github.com/nazar256/go-amqp-sniffer/broker"
	"github.com/streadway/amqp"
)

const (
	messageID    = "message_id"
	appID        = "test_app_id"
	header1Value = "h1_value"
	header2Value = "h2_value"
	contentType  = "application/json"
)

type ListenerMock struct {
	IsCancelled bool

	headers map[string]interface{}
	time    time.Time
	ch      chan amqp.Delivery
}

func newListenerMock(timestamp time.Time) *ListenerMock {
	return &ListenerMock{
		headers: map[string]interface{}{
			"header_1": header1Value,
			"header_2": header2Value,
		},
		time: timestamp,
		ch:   make(chan amqp.Delivery),
	}
}

func (l ListenerMock) Listen() broker.AMQPMsgChan {
	return l.ch
}

func (l ListenerMock) Cancel() error {
	close(l.ch)
	l.IsCancelled = true

	return nil
}

func (l ListenerMock) enqueue(msgBody string) {
	l.ch <- amqp.Delivery{
		MessageId:       messageID,
		AppId:           appID,
		Timestamp:       l.time,
		Headers:         l.headers,
		ContentType:     contentType,
		ContentEncoding: "",

		Body: []byte(msgBody),
	}
}
