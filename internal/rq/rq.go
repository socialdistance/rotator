package rq

import (
	"context"
	"fmt"
	rq "github.com/rabbitmq/amqp091-go"
	"rotator/internal/app"
)

type Rabbit struct {
	exchange string
	queue    string
	consumer string
	channel  *rq.Channel
	logger   app.Logger
}

func NewRabbit(ctx context.Context, url, exchange, queue string, logger app.Logger) (*Rabbit, error) {
	conn, err := rq.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ on %s: %w", url, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open RabbitMQ Channel on %s: %w", url, err)
	}

	if len(exchange) > 0 {
		err = ch.ExchangeDeclare(
			exchange,
			rq.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare an exchange %s: %w", exchange, err)
		}
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue %s: %w", queue, err)
	}

	err = ch.QueueBind(
		q.Name,
		q.Name,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	go func() {
		<-ctx.Done()
		ch.Close()
		conn.Close()
	}()

	return &Rabbit{
		exchange: exchange,
		queue:    queue,
		consumer: "rotator-consumer",
		channel:  &rq.Channel{},
		logger:   logger,
	}, nil
}
