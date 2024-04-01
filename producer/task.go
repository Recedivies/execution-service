package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/exception"
)

type Task struct {
	ch *amqp.Channel
}

// NewTask instantiates the Task repository.
func NewTask(channel *amqp.Channel) (*Task, error) {
	return &Task{
		ch: channel,
	}, nil
}

func (t *Task) Publish(ctx context.Context, exchange, routingKey string, e ...interface{}) error {
	var b bytes.Buffer

	if e != nil {
		jsonData, err := json.Marshal(e)
		if err != nil {
			return err
		}

		if _, err := b.Write(jsonData); err != nil {
			return err
		}
	}

	var deliveryMode uint8
	if routingKey == "dlq.task" {
		deliveryMode = amqp.Persistent
	} else {
		deliveryMode = amqp.Transient
	}

	err := t.ch.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: deliveryMode,
			ContentType:  "application/x-encoding-gob",
			Body:         b.Bytes(),
			Timestamp:    time.Now(),
		})
	if err != nil {
		return exception.WrapErrorf(err, exception.ErrorCodeUnknown, "ch.Publish")
	}

	return nil
}
