package rabbitmq

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ(ctx context.Context, url string) (*amqp091.Connection, *amqp091.Channel, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	_, err = channel.QueueDeclare("wa_otp_queue", true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	return conn, channel, nil
}

func PublishMessage(ctx context.Context, ch *amqp091.Channel, queueName string, body []byte) error {
	if err := ch.PublishWithContext(ctx, "otp", queueName, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body: body,
	}); err != nil {
		return err
	}

	return nil
} 