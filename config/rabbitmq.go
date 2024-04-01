package config

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/exception"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(rabbitMQServerURL string) *RabbitMQ {
	conn, err := amqp.Dial(rabbitMQServerURL)
	exception.FatalIfNeeded(err, "amqp.Dial")

	ch, err := conn.Channel()
	exception.FatalIfNeeded(err, "conn.Channel")

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
	}
}

func (r *RabbitMQ) Close() {
	r.Connection.Close()
}
