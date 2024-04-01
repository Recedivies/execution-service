package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/config"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/mail"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/model"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/producer"
	"gorm.io/gorm"
)

type TaskProcessor interface {
	ListenAndServe() error
	SendCallback(model.Job) error
	Shuwdown(ctx context.Context) error
}

type RabbitMQTaskProcessor struct {
	DB        *gorm.DB
	mail      mail.EmailSender
	rmq       *config.RabbitMQ
	msgBroker *producer.Task
	done      chan struct{}
}

func NewRabbitMQTaskProcessor(mail mail.EmailSender, DB *gorm.DB, rmq *config.RabbitMQ, msgBroker *producer.Task) TaskProcessor {
	return &RabbitMQTaskProcessor{
		DB:        DB,
		mail:      mail,
		rmq:       rmq,
		msgBroker: msgBroker,
		done:      make(chan struct{}),
	}
}

func (s *RabbitMQTaskProcessor) SendCallback(job model.Job) error {
	callbackURL := "https://email.requestcatcher.com/test"

	payload := map[string]interface{}{
		"status":  "success",
		"message": "Task processed successfully",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Error sending callback")
		return err
	}

	_, err = http.Post(callbackURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Error().Err(err).Msg("Error sending callback")
		return err
	}

	log.Info().Msg("Callback sent successfully")
	return nil
}

func (s *RabbitMQTaskProcessor) Shuwdown(ctx context.Context) error {
	log.Info().Msg("Shutting down server")

	s.rmq.Channel.Cancel(queueName, false)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context.Done: %w", ctx.Err())

		case <-s.done:
			return nil
		}
	}
}
