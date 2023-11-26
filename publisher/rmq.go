package publisher

import (
	"context"
	"log"
	"time"

	"github.com/furdarius/rabbitroutine"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PublisherTypeRabbitMQ = "rabbitmq"
)

const (
	RabbitMQCompressionHeader      = "compression"
	RabbitMQCompressionZstd        = "application/zstd"
	RabbitMQContentType            = "application/json"
	RabbitMQReconnectAttempts      = 5               // How many times to attempt reconnecting
	RabbitMQReconnectWait          = 1 * time.Second // How long to wait before another reconnect attempt
	RabbitMQPublishBaseDelay       = 2 * time.Second // Base duration for every retry
	RabbitMQMaxPublishRetries uint = 5               // Max number of retries
)

type RabbitMQOpts struct {
	URL        string
	Exchange   string
	RoutingKey string
	Compress   bool
}

type RabbitMQ struct {
	publisher  rabbitroutine.Publisher
	exchange   string
	routingKey string
	compress   bool
}

var connector = rabbitroutine.NewConnector(rabbitroutine.Config{
	Wait:              RabbitMQReconnectWait,     // how long wait between reconnect
	ReconnectAttempts: RabbitMQReconnectAttempts, // max attempts for dialling
})

var rabbitroutinePublisher = rabbitroutine.NewRetryPublisher(
	rabbitroutine.NewEnsurePublisher(rabbitroutine.NewPool(connector)),
	rabbitroutine.PublishMaxAttemptsSetup(RabbitMQMaxPublishRetries),
	rabbitroutine.PublishDelaySetup(
		rabbitroutine.LinearDelay(RabbitMQPublishBaseDelay), // time.Duration(attempt) * delay
	),
)

func NewRabbitMQPublisher(ctx context.Context, opts *RabbitMQOpts) Publisher {
	go connector.Dial(ctx, opts.URL)
	return &RabbitMQ{
		publisher:  rabbitroutinePublisher,
		exchange:   opts.Exchange,
		routingKey: opts.RoutingKey,
		compress:   opts.Compress,
	}
}

func (r *RabbitMQ) Publish(ctx context.Context, payload Payload) error {
	body, err := payload.Bytes()
	if err != nil {
		log.Println("error while compressing payload before publishing to RabbitMQ", err)
		return err
	}

	message := amqp.Publishing{
		DeliveryMode: amqp.Transient,
		ContentType:  RabbitMQContentType,
		Body:         body,
	}

	if r.compress {
		message.Headers = amqp.Table{
			RabbitMQCompressionHeader: RabbitMQCompressionZstd,
		}
	}

	if err := r.publisher.Publish(ctx,
		r.exchange,   // Exchange
		r.routingKey, // Routing key
		message,
	); err != nil {
		log.Println("error while publishing to RabbitMQ", err)
		return err
	}
	log.Println("published to RabbitMQ")
	return nil
}
