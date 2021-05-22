package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nazar256/go-amqp-sniffer/broker"
	"github.com/nazar256/go-amqp-sniffer/flow"
	"github.com/nazar256/go-amqp-sniffer/sniffer"
	"github.com/spf13/cobra"
)

const (
	flagURL             = "url"
	flagExchange        = "exchange"
	flagRoutingKey      = "routing-key"
	flagListenQueue     = "listen-queue"
	flagPrefetch        = "prefetch"
	flagFormat          = "format"
	flagParseJSON       = "parse-json"
	flagStopOnErr       = "stop-on-error"
	flagPersistentQueue = "persistent-queue"
)

func initRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   `go-amqp-sniffer [options]`,
		Args:  cobra.NoArgs,
		Short: "Listens to subscribed messages in AMQP and outputs them in selected format",
		Long: `Creates a temporary queue for sniffing, binds it to specified exchange with specified routing keys.
To listen all the events easily an exchange must be of type "topic" (routing key is "#") or fanout (empty routing key).
Otherwise (type "direct"), you have to specify each routing key to listen all the events.
The most flexible format is JSON. It's compatible with mongoimport`,
		RunE: runRoot,
	}

	rootCmd.Flags().String(
		flagURL,
		"amqp://guest:guest@localhost:5672/",
		"AMQP url with credentials, host, port and vhost",
	)

	rootCmd.Flags().String(
		flagExchange,
		"",
		"Exchange to listen events from",
	)

	err := rootCmd.MarkFlagRequired(flagExchange)
	flow.FailOnError(err, "failed to mark exchange flag as required")

	rootCmd.Flags().StringSlice(
		flagRoutingKey,
		[]string{"#"},
		"Routing keys to subscribe. Default is # which subscribes for all events on topic exchange",
	)

	rootCmd.Flags().String(
		flagListenQueue,
		"amqp-sniffer",
		"queue name for receiving events (messages)",
	)

	rootCmd.Flags().Int(
		flagPrefetch,
		10,
		"prefetch count for AMQP",
	)

	rootCmd.Flags().String(
		flagFormat,
		"json",
		"output (stdout) format (json, csv)",
	)

	rootCmd.Flags().Bool(
		flagParseJSON,
		false,
		"parse content body as JSON, does not have effect with --format csv",
	)

	rootCmd.Flags().Bool(
		flagStopOnErr,
		false,
		"exit program even when non-critical error occurs (i.e. parsing or serialization error)",
	)

	rootCmd.Flags().Bool(
		flagPersistentQueue,
		false,
		"do not delete queue automatically when program exits (normally or on flow)",
	)

	return rootCmd
}

func runRoot(cmd *cobra.Command, _ []string) error {
	var format sniffer.Format

	specifiedFormat, _ := cmd.Flags().GetString(flagFormat)
	switch specifiedFormat {
	case "json":
		format = sniffer.JSON
	case "csv":
		format = sniffer.CSV
	default:
		log.Fatalf("Specified output specifiedFormat %s is invalid, valid formats are json or csv", specifiedFormat)
	}

	url, _ := cmd.Flags().GetString(flagURL)
	exchange, _ := cmd.Flags().GetString(flagExchange)
	routingKeys, _ := cmd.Flags().GetStringSlice(flagRoutingKey)
	listenQueueName, _ := cmd.Flags().GetString(flagListenQueue)
	persistentQueue, _ := cmd.Flags().GetBool(flagPersistentQueue)
	prefetch, _ := cmd.Flags().GetInt(flagPrefetch)
	listener := broker.NewListener(&broker.Config{
		URL:             url,
		ExchangeName:    exchange,
		RoutingKeys:     routingKeys,
		QueueName:       listenQueueName,
		PersistentQueue: persistentQueue,
		Prefetch:        prefetch,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-termChan
		cancel()
	}()

	wg := &sync.WaitGroup{}

	stopOnError, _ := cmd.Flags().GetBool(flagStopOnErr)
	parseJSON, _ := cmd.Flags().GetBool(flagParseJSON)

	sniffer.Sniff(ctx, wg, listener, os.Stdout, sniffer.Config{
		OutputFormat: format,
		StopOnErr:    stopOnError,
		Parse:        parseJSON,
	})

	wg.Wait()

	log.Println("Sniffer is stopped.")

	return nil
}
