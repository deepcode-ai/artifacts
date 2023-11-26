package publisher

import (
	"context"
	"errors"
	"reflect"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

const RESULT_RMQ_URL = "amqp://localhost:5672/"

type MockAMQPPublisher struct {
	err error
}

func (p *MockAMQPPublisher) Publish(_ context.Context, _, _ string, _ amqp.Publishing) error {
	return p.err
}

func initializeRabbitMQ() error {
	conn, err := amqp.Dial(RESULT_RMQ_URL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare exchange `celery` to listen for results pushed from rome
	if err := ch.ExchangeDeclare("celery", "direct", true, false, false, false, nil); err != nil {
		return err
	}

	// Declare `celery` queue
	if _, err := ch.QueueDeclare("celery", true, false, false, false, nil); err != nil {
		return err
	}

	// Bind `celery` queue to `celery` exchange
	if err := ch.QueueBind("celery", "celery", "celery", false, nil); err != nil {
		return err
	}
	return nil
}

func TestNewRabbitMQPublisher(t *testing.T) {
	if err := initializeRabbitMQ(); err != nil {
		t.Fatalf("failed to initialize RabbitMQ: %v", err)
	}
	ctx := context.Background()

	v := NewRabbitMQPublisher(ctx, &RabbitMQOpts{
		URL:        RESULT_RMQ_URL,
		Exchange:   "celery",
		RoutingKey: "celery",
	})
	p, ok := v.(*RabbitMQ)
	if p == nil || !ok {
		t.Fatal("NewRabbitMQPublisher() did not return a RabbitMQ publisher")
	}

	if p.exchange != "celery" {
		t.Errorf("NewRabbitMQPublisher() exchange = %v, want %v", p.exchange, "celery")
	}
	if p.routingKey != "celery" {
		t.Errorf("NewRabbitMQPublisher() routingKey = %v, want %v", p.routingKey, "celery")
	}

	if !reflect.DeepEqual(p.publisher, rabbitroutinePublisher) {
		t.Errorf("NewRabbitMQPublisher() publisher = %v, want %v", p.publisher, rabbitroutinePublisher)
	}
}

func TestRabbitMQ_Publish(t *testing.T) {
	type fields struct {
		publisher Publisher
	}

	type args struct {
		ctx     context.Context
		payload Payload
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				publisher: &RabbitMQ{
					publisher: &MockAMQPPublisher{},
				},
			},
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					payload: []byte("test"),
				},
			},
		},
		{
			name: "uncompressible payload",
			fields: fields{
				publisher: &RabbitMQ{
					publisher: &MockAMQPPublisher{},
				},
			},
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					err: errors.New("test-error"),
				},
			},
			wantErr: true,
		},
		{
			name: "error publishing",
			fields: fields{
				publisher: &RabbitMQ{
					publisher: &MockAMQPPublisher{
						err: errors.New("test-error"),
					},
				},
			},
			args: args{
				ctx: context.Background(),
				payload: &MockPayload{
					payload: []byte("test"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields.publisher
			if err := r.Publish(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("RabbitMQ.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
