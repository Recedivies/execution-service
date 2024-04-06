package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/config"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/mail"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/model"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/producer"
	"gorm.io/gorm"
)

type TaskProcessor interface {
	ListenTaskRetryAndServe() error
	ListenTaskAndServe() error
	SendCallback(job model.Job, taskExecutionHistory model.TaskExecutionHistory, err error) error
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

func (s *RabbitMQTaskProcessor) SendCallback(job model.Job, taskExecutionHistory model.TaskExecutionHistory, err error) error {
	executionTime := time.Now()
	if taskExecutionHistory != (model.TaskExecutionHistory{}) {
		executionTime = taskExecutionHistory.ExecutionTime
	}

	status := "success"
	message := "task processed successfully"
	if err != nil {
		status = "failed"
		message = "task processing failed"
	}

	payload := map[string]interface{}{
		"status":         status,
		"message":        message,
		"execution_time": executionTime.Format(time.RFC3339), // Convert time to string format
	}

	if err != nil {
		payload["error"] = err.Error()
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("json.Marshal")
		return err
	}

	log.Info().Str("payload", string(jsonPayload)).Str("URL", job.CallbackUrl).Msg("Sending to callback URL")

	_, err = http.Post(job.CallbackUrl+job.Id, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Error().Err(err).Msg("Error sending callback")
		return err
	}

	log.Info().Msg("Callback sent successfully")
	return nil
}

func (s *RabbitMQTaskProcessor) Shuwdown(ctx context.Context) error {
	log.Info().Msg("Shutting down server")

	s.rmq.Channel.Cancel(queueTask, false)
	s.rmq.Channel.Cancel(queueTaskRetry, false)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context.Done: %w", ctx.Err())

		case <-s.done:
			return nil
		}
	}
}
