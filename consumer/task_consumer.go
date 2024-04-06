package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/model"
)

const (
	queueTask    = "q.task"
	consumerTask = "task"

	FAILED_STATUS    = "FAILED"
	COMPLETED_STATUS = "COMPLETED"
)

func (processor *RabbitMQTaskProcessor) ListenTaskAndServe() error {
	msgs, err := processor.rmq.Channel.Consume(
		queueTask,    // queue
		consumerTask, // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("channel.Consume %w", err)
	}

	ctx := context.Background()

	go func() {
		for msg := range msgs {
			log.Info().Msg(fmt.Sprintf("Received message: %s", msg.RoutingKey))

			var err error
			var nack bool = false
			var job model.Job
			var taskExecutionHistory model.TaskExecutionHistory

			var jobId string = string(msg.Body)
			log.Info().Str("Job Id", jobId).Msg("processed task")

			switch msg.RoutingKey {
			case "task.SEND_EMAIL":
				result := processor.DB.Where("id = ?", jobId).Find(&job)
				if result.Error != nil {
					err = result.Error
					log.Error().Err(result.Error).Msg("DB.Model")
					nack = true
					break
				}

				result = processor.DB.First(&taskExecutionHistory, "job_id = ?", jobId)
				if result.Error != nil {
					err = result.Error
					log.Error().Err(result.Error).Msg("DB.Model")
					nack = true
					break
				}

				var configMap map[string]interface{}
				err = json.Unmarshal([]byte(job.Config), &configMap)
				if err != nil {
					log.Error().Err(err).Msg("Failed to unmarshal job.Config JSON")
					nack = true
					break
				}

				subject := "Welcome to Job Scheduler Service"
				content := configMap["message"].(string)
				to := []string{configMap["receiver"].(string)}

				err = processor.mail.SendEmail(subject, content, to, nil, nil, nil)
				if err != nil {
					log.Error().Err(err).Msg("mail.SendEmail")
					nack = true
					break
				}

				taskExecutionHistory.ExecutionTime = time.Now()
				taskExecutionHistory.Status = COMPLETED_STATUS
				processor.DB.Save(&taskExecutionHistory)

				log.Info().Msg("Successful completion of task execution")
				processor.SendCallback(job, taskExecutionHistory, nil)

			default:
				nack = true
			}

			if nack {
				log.Info().Msg("NAcking")
				if taskExecutionHistory == (model.TaskExecutionHistory{}) || job == (model.Job{}) {
					processor.SendCallback(job, taskExecutionHistory, errors.New("empty model"))
					msg.Nack(false, false)
					continue
				}

				if taskExecutionHistory.RetryCount <= job.MaxRetryCount {
					taskExecutionHistory.RetryCount += 1
					err := processor.msgBroker.Publish(ctx, "x.task", "task-retry."+job.JobType, jobId)
					if err != nil {
						log.Error().Err(err).Msg("msgBroker.Publish")
					}
					log.Info().Msg("Published to task retry queue")
				} else {
					// produce a Job to dlx.task
					taskExecutionHistory.Status = FAILED_STATUS
					err := processor.msgBroker.Publish(ctx, "dlx.task", "dlq.task", job)
					if err != nil {
						log.Error().Err(err).Msg("msgBroker.Publish")
					}
				}
				taskExecutionHistory.ExecutionTime = time.Now()
				processor.DB.Save(&taskExecutionHistory)
				processor.SendCallback(job, taskExecutionHistory, err)
				msg.Nack(false, false)
			} else {
				log.Info().Msg("Acking")
				msg.Ack(false)
			}
		}

		processor.done <- struct{}{}
	}()

	return nil
}
