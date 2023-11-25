package artifacts

import (
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/streadway/amqp"
)

var ConsumerChannelActive = true

// SetupRMQConnection sets up the RabbitMQ connection and channels. It also spawns a goroutine that listens for any
// "Close" events from the broker.
func SetupRMQConnection(retryFunc func() error, brokerUrl, exchangeName, exchangeType string) (*amqp.Connection, *amqp.Channel, error) {
	// Establish connection with RabbitMQ
	rmqConn, err := dialRMQ(brokerUrl)
	if err != nil {
		return nil, nil, err
	}

	// Spawn a goroutine to handle any RabbitMQ related errors.
	go handleRabbitMQErrors(rmqConn, retryFunc)

	// Create Channel
	rmqChannel, err := createChannel(rmqConn)
	if err != nil {
		return nil, nil, err
	}

	return rmqConn, rmqChannel, declareRealtimeExchange(rmqChannel, exchangeName, exchangeType)
}

// createChannel opens channel with janus virtualhost connection.
func createChannel(rmqConn *amqp.Connection) (*amqp.Channel, error) {
	rmqChan, err := rmqConn.Channel()
	return rmqChan, err
}

// Declare exchange `janus-realtime-notifications` to listen for realtime notification pushes from asgard
func declareRealtimeExchange(rmqChan *amqp.Channel, exchangeName, exchangeType string) (err error) {
	if err = rmqChan.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Error declaring exchange: janus-realtime-notifications on vhost: janus. Error: %s", err)
	}
	return nil
}

// handleRabbitMQErrors looks out for any RabbitMQ errors/closure and re-establishes connection
// and initializes the realtime consumer.
func handleRabbitMQErrors(rmqConn *amqp.Connection, consumerInitFn func() error) {
	rmqError := <-rmqConn.NotifyClose(make(chan *amqp.Error))
	rmqCloseErr := fmt.Errorf("RabbitMQ connection closed")

	// RabbitMQ connection closure triggered. Log and send a sentry alert.
	if rmqError != nil {
		rmqCloseErr = fmt.Errorf("RMQ Connection closed. Code: %d: Reason: %s", rmqError.Code, rmqError.Reason)
	}
	log.Println(rmqCloseErr)
	sentry.CaptureException(rmqCloseErr)

	// Check if the error is a "ConnectionError" or a "ChannelError".
	// Other error codes don't trigger retrial.
	if !ConsumerChannelActive || rmqError != nil && isBrokerError(rmqError.Code) {
		for {
			if retryError := retryBrokerConnections(consumerInitFn); retryError == nil {
				break
			}
		}
	}
}

// retryBrokerConnections retries broker connection function once the broker error/closure is detected.
func retryBrokerConnections(consumerSetupFn func() error) error {
	// Setup consumers
	if reloadErr := consumerSetupFn(); reloadErr != nil {
		log.Println("Failed to restore RMQ consumers: ", reloadErr)
		// raven.CaptureErrorAndWait(fmt.Errorf("Failed to restore RMQ consumers: %v", reloadErr), nil)
		return reloadErr
	}
	return nil
}

// isBrokerError checks if the error code is a signifact broker error to re-establish connection and channel.
func isBrokerError(code int) bool {
	switch code {
	case
		// Channel Errors
		311, // amqp.ContentTooLarge
		313, // amqp.NoConsumers
		403, // amqp.AccessRefused
		404, // amqp.NotFound
		405, // amqp.ResourceLocked
		406: // amqp.PreconditionFailed
		return true

	case
		// Connection Errors
		320, // amqp.ConnectionForced
		402, // amqp.InvalidPath
		501, // amqp.FrameError
		502, // amqp.SyntaxError
		503, // amqp.CommandInvalid
		504, // amqp.ChannelError
		505, // amqp.UnexpectedFrame
		506, // amqp.ResourceError
		530, // amqp.NotAllowed
		540, // amqp.NotImplemented
		541: // amqp.InternalError
		return true
	default:
		return false
	}
}
