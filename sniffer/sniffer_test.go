package sniffer_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/nazar256/go-amqp-sniffer/sniffer"
)

func TestSnifferReceivesMessagesAndLogsThem(t *testing.T) {
	type caseConfig struct {
		cfg         sniffer.Config
		msgBodies   []string
		outputLines []string
	}

	type testTable map[string]caseConfig

	tests := testTable{
		"single message is sniffed": {
			sniffer.Config{OutputFormat: sniffer.JSON, Parse: true},
			[]string{"{\"root\":{\"something\":12345}}"},
			[]string{"{\"MessageID\":\"%s\",\"AppID\":\"%s\",\"Timestamp\":\"%s\",\"Headers\":{\"header_1\":\"%s\"," +
				"\"header_2\":\"%s\"},\"Payload\":\"\",\"ParsedPayload\":{\"root\":{\"something\":12345}}," +
				"\"ContentType\":\"application/json\",\"ContentEncoding\":\"\"}\n"},
		},
		"multiple messages are sniffed": {
			sniffer.Config{OutputFormat: sniffer.JSON, Parse: true},
			[]string{
				"{\"root\":{\"something\":12345}}",
				"{\"single_field\":\"value\"}",
			},
			[]string{
				"{\"MessageID\":\"%s\",\"AppID\":\"%s\",\"Timestamp\":\"%s\"," +
					"\"Headers\":{\"header_1\":\"%s\",\"header_2\":\"%s\"}," +
					"\"Payload\":\"\",\"ParsedPayload\":{\"root\":{\"something\":12345}}," +
					"\"ContentType\":\"application/json\",\"ContentEncoding\":\"\"}\n",
				"{\"MessageID\":\"%s\",\"AppID\":\"%s\",\"Timestamp\":\"%s\"," +
					"\"Headers\":{\"header_1\":\"%s\",\"header_2\":\"%s\"}," +
					"\"Payload\":\"\",\"ParsedPayload\":{\"single_field\":\"value\"}," +
					"\"ContentType\":\"application/json\",\"ContentEncoding\":\"\"}\n",
			},
		},
		"disabled parsing fill original payload": {
			sniffer.Config{OutputFormat: sniffer.JSON, Parse: false},
			[]string{"{\"document\":{\"value\":1.234}}"},
			[]string{"{\"MessageID\":\"%s\",\"AppID\":\"%s\",\"Timestamp\":\"%s\"," +
				"\"Headers\":{\"header_1\":\"%s\",\"header_2\":\"%s\"}," +
				"\"Payload\":\"{\\\"document\\\":{\\\"value\\\":1.234}}\",\"ParsedPayload\":null," +
				"\"ContentType\":\"application/json\",\"ContentEncoding\":\"\"}\n"},
		},
		"parsing failure leaves original payload": {
			sniffer.Config{OutputFormat: sniffer.JSON, Parse: true},
			[]string{"{\"invalid_json:123"},
			[]string{"{\"MessageID\":\"%s\",\"AppID\":\"%s\",\"Timestamp\":\"%s\"," +
				"\"Headers\":{\"header_1\":\"%s\",\"header_2\":\"%s\"}," +
				"\"Payload\":\"{\\\"invalid_json:123\",\"ParsedPayload\":null," +
				"\"ContentType\":\"application/json\",\"ContentEncoding\":\"\"}\n"},
		},
		"csv": {
			sniffer.Config{OutputFormat: sniffer.CSV},
			[]string{"{\"key\":\"value\"}"},
			[]string{"%s,%s,%s,\"{\"\"header_1\"\":\"\"%s\"\",\"\"header_2\"\":\"\"%s\"\"}\"," +
				"\"{\"\"key\"\":\"\"value\"\"}\",application/json,\n"},
		},
	}

	for caseName, test := range tests {
		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}
		output := &bytes.Buffer{}
		now := time.Now()
		listener := newListenerMock(now)

		go func(test caseConfig) {
			for _, msgBody := range test.msgBodies {
				listener.enqueue(msgBody)
			}

			cancel()
		}(test)

		sniffer.Sniff(ctx, wg, listener, output, test.cfg)

		var expectedOutput string
		for _, outputPattern := range test.outputLines {
			expectedOutput += fmt.Sprintf(
				outputPattern,
				messageID,
				appID,
				now.Format(time.RFC3339Nano),
				header1Value,
				header2Value,
			)
		}

		if output.String() != expectedOutput {
			log.Fatalf("For case \"%s\" received output %s but expected %s", caseName, output, expectedOutput)
		}
	}
}
