package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/exception"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/model"
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
	var b []byte

	switch exchange {
	case "x.task":
		if len(e) == 0 {
			return fmt.Errorf("no string provided for 'x.task' exchange")
		}
		str, ok := e[0].(string)
		if !ok {
			return fmt.Errorf("expected string for 'x.task' exchange")
		}
		b = []byte(str)
	case "dlx.task":
		if len(e) == 0 {
			return fmt.Errorf("no model provided for 'dlx.task' exchange")
		}
		model, ok := e[0].(model.Job)
		if !ok {
			return fmt.Errorf("expected model struct for 'dlx.task' exchange")
		}

		jsonData, err := json.Marshal(model)
		if err != nil {
			return err
		}
		b = append(b, jsonData...)
	}

	var deliveryMode uint8 = amqp.Transient
	if routingKey == "dlq.task" {
		deliveryMode = amqp.Persistent
	}

	log.Info().Str("body message", string(b)).Msg("Published data")

	err := t.ch.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: deliveryMode,
			ContentType:  "application/x-encoding-gob",
			Body:         b,
			Timestamp:    time.Now(),
		})
	if err != nil {
		return exception.WrapErrorf(err, exception.ErrorCodeUnknown, "ch.Publish")
	}

	return nil
}
