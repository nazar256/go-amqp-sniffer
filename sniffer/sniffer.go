package sniffer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/nazar256/go-amqp-sniffer/broker"
	"github.com/nazar256/go-amqp-sniffer/flow"
	"github.com/nazar256/parapipe"
	"github.com/streadway/amqp"
)

type Format int

const (
	JSON Format = iota
	CSV
)

type Config struct {
	OutputFormat Format
	Parse        bool
	StopOnErr    bool
}

type OutputRecord struct {
	MessageID       string
	AppID           string
	Timestamp       time.Time
	Headers         map[string]interface{}
	Payload         string
	ParsedPayload   map[string]interface{}
	ContentType     string
	ContentEncoding string
}

type Listener interface {
	Listen() broker.AMQPMsgChan
	Cancel() error
}

func Sniff(ctx context.Context, wg *sync.WaitGroup, listener Listener, out io.Writer, cfg Config) {
	serialize := getSerializer(cfg.OutputFormat)
	pipeline := parapipe.NewPipeline(parapipe.Config{Concurrency: runtime.NumCPU()})

	// cancel amqp consumer in case of canceled context (i.e. user pressed Ctrl+C)
	go func() {
		<-ctx.Done()

		err := listener.Cancel()
		flow.FailOnError(err, "failed to cancel consumer")
	}()

	// defining the pipeline stages
	pipeline.

		// prepares record structures, including optional message body parsing
		Pipe(func(msg interface{}) interface{} {
			amqpMsg := msg.(amqp.Delivery)

			record := OutputRecord{
				MessageID:       amqpMsg.MessageId,
				AppID:           amqpMsg.AppId,
				Timestamp:       amqpMsg.Timestamp,
				Headers:         amqpMsg.Headers,
				ContentType:     amqpMsg.ContentType,
				ContentEncoding: amqpMsg.ContentEncoding,
			}

			var unmarshalErr error
			if cfg.Parse {
				unmarshalErr = json.Unmarshal(amqpMsg.Body, &record.ParsedPayload)
				if unmarshalErr != nil && cfg.StopOnErr {
					log.Fatalf("Failed to parse body %s, %s", amqpMsg.Body, unmarshalErr)
				}
			}
			if !cfg.Parse || cfg.OutputFormat == CSV || unmarshalErr != nil {
				record.Payload = string(amqpMsg.Body)
			}

			return record
		}).

		// serializes each record in batches to the specified format
		Pipe(func(msg interface{}) interface{} {
			record := msg.(OutputRecord)
			line, err := serialize(&record)
			if err != nil && cfg.StopOnErr {
				log.Fatalf("Failed to serialize record %s, %s", record, err)
			}

			return line
		})

	bindListenerToPipeline(listener, wg, pipeline)

	// executes pipeline until its input is closed)
	for line := range pipeline.Out() {
		l := line.([]byte)
		_, err := out.Write(l)
		flow.FailOnError(err, "Failed to write to output")
		wg.Done()
	}
}

func bindListenerToPipeline(listener Listener, wg *sync.WaitGroup, pipeline *parapipe.Pipeline) {
	go func() {
		for amqpMsg := range listener.Listen() {
			wg.Add(1)
			pipeline.Push(amqpMsg)
		}

		pipeline.Close()
	}()
}
