package datastore

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

var Rabbit *RabbitMQ

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func InitRabbitMQ() error {
	rabbitURL := os.Getenv("RABBIT_URL")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	Rabbit = &RabbitMQ{
		conn:    conn,
		channel: channel,
	}

	return nil
}

func (r *RabbitMQ) Close() error {
	err := r.channel.Close()
	if err != nil {
		return err
	}
	err = r.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, queueName, message string) error {
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	err := r.channel.PublishWithContext(ctx, "", queueName, false, false, publishing)
	if err != nil {
		return err
	}

	return nil
}

func SendNationalToRabbit(national string) error {
	queueName := "national"
	ctx := context.Background()
	err := Rabbit.PublishMessage(ctx, queueName, national)
	if err != nil {
		return err
	}

	return nil
}
