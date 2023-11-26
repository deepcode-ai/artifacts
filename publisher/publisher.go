// Publisher is a message publisher that publishes resutls to Asgard.
package publisher

import "context"

type Payload interface {
	Bytes() ([]byte, error)
}

type Publisher interface {
	Publish(ctx context.Context, payload Payload) error
}
