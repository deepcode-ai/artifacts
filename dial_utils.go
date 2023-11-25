package artifacts

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

var maxRetryCount = 12

// dialRMQ creates a connection with RabbiMQ and stores the amqp connection
// object for usage by other functions
func dialRMQ(brokerUrl string) (rmqConn *amqp.Connection, err error) {
	retryTimeout := 1
	retryCount := 0
	var retryDuration time.Duration

	for {
		if rmqConn, err = amqp.Dial(brokerUrl); err == nil {
			log.Println("Successfully reconnected to RabbitMQ")
			return rmqConn, nil
		}

		log.Printf("Failed to connect to RMQ, Error: %v", err)
		if retryCount > maxRetryCount {
			return nil, err
		}
		// Wait for retrying in time intervals based on fibonacci series
		retryTimeout, retryDuration = GetRetryTimeout(retryTimeout)
		log.Printf("Retrying connection in %d seconds", retryTimeout)
		time.Sleep(retryDuration)
		retryCount++
	}
}
