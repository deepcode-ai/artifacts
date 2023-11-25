package consumer

import (
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

// Single function for all the consumer related operations
func Consume(rmqConn *amqp.Connection, queueName, routingKey string, block chan error, processMessage func(amqp.Delivery) error) {
	rmqChannel, err := rmqConn.Channel()
	if err != nil {
		sentry.CaptureException(err)
		block <- fmt.Errorf("error opening channel: %s", err)
	}

	defer rmqChannel.Close()

	// Declare analysis run queue
	queueDeclare, err := rmqChannel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when used
		false,     // Exclusive
		false,     // No wait
		nil,       // Arguments
	)
	if err != nil {
		sentry.CaptureException(err)
		block <- fmt.Errorf("error declaring queue: %s %s", queueName, err)
	}

	// Bind analysis run queue to atlas-jobs exchange
	err = rmqChannel.QueueBind(
		queueName,                           // Queue name
		routingKey,                          // Routing key
		viper.GetString("app.rmq.exchange"), // Exchange name
		false,                               // No wait
		nil,                                 // Arguments
	)
	if err != nil {
		sentry.CaptureException(err)
		block <- fmt.Errorf("error binding queue: %s %s", queueName, err)
	}

	// Listen for messages to consume from analysis-run queue
	messages, err := rmqChannel.Consume(
		queueDeclare.Name, // Queue name
		"",                // Consumer
		true,              // Auto ackowledge
		false,             // Exclusive
		false,             // No local
		false,             // No wait
		nil,               // Arguments
	)
	if err != nil {
		sentry.CaptureException(err)
		block <- fmt.Errorf("error registering consumer for queue: %s %s", queueName, err)
	}

	go func() {
		for message := range messages {
			if err := processMessage(message); err != nil {
				log.Println(err)
				sentry.CaptureException(err)
				continue
			}
		}
	}()

	log.Printf("Consumers initialized for queue: %s", queueName)

	<-block
}
