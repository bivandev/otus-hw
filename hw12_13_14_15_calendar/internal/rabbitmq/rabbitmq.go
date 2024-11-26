package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Close()
	DeclareQueue() error
	Publish(body []byte) error
	Consume() (<-chan []byte, error)
}

const queueName = "event_notifications"

type RabbitMQ struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewQueue(url string) (Queue, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{conn: conn, channel: channel}, nil
}

func (r *RabbitMQ) DeclareQueue() error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (r *RabbitMQ) Publish(body []byte) error {
	return r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Consume() (<-chan []byte, error) {
	msgs, err := r.channel.Consume(
		queueName,
		"",    // Имя потребителя
		true,  // Автоматическое подтверждение
		false, // Эксклюзивный доступ
		false, // Локальные сообщения (обычно false)
		false, // Ожидание
		nil,   // Дополнительные аргументы
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages from queue %s: %w", queueName, err)
	}

	out := make(chan []byte)

	go func() {
		defer close(out)
		for msg := range msgs {
			out <- msg.Body
		}
	}()

	return out, nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
