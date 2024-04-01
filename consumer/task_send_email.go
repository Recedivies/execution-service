package consumer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/model"
)

const (
	queueName    = "q.task"
	consumerName = "task"

	FAILED_STATUS    = "FAILED"
	COMPLETED_STATUS = "COMPLETED"
)

func (processor *RabbitMQTaskProcessor) ListenAndServe() error {
	msgs, err := processor.rmq.Channel.Consume(
		queueName,    // queue
		consumerName, // consumer
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

	cnt := 0

	go func() {
		for msg := range msgs {
			cnt += 1
			log.Info().Msg(fmt.Sprintf("Received message: %s", msg.RoutingKey))

			var nack bool
			var job model.Job
			var taskExecutionHistory model.TaskExecutionHistory

			// send email
			subject := "Welcome to Job Scheduler Service"
			content := `Hello,<br/>
						Thank you for registering with us!<br/>
						`
			to := []string{"ahmadhihastiputra@gmail.com"}

			var jobId string = string(msg.Body)
			log.Info().Str("Job Id", jobId).Msg("processed task")

			switch msg.RoutingKey {
			case "tasks":
				result := processor.DB.Model(model.Job{Id: jobId}).First(&job)
				if result.Error != nil {
					log.Error().Err(result.Error).Msg("")
					nack = true
					break
				}

				result = processor.DB.Model(model.TaskExecutionHistory{JobId: jobId}).First(&taskExecutionHistory)
				if result.Error != nil {
					log.Error().Err(result.Error).Msg("")
					nack = true
					break
				}

				err = processor.mail.SendEmail(subject, content, to, nil, nil, nil)
				if err != nil {
					fmt.Println(err.Error())
					log.Error().Err(err).Msg("")
					nack = true
					break
				}

				taskExecutionHistory.Status = COMPLETED_STATUS
				processor.DB.Save(&taskExecutionHistory)

				log.Info().Msg("Successful completion of task execution")
				processor.SendCallback(job)

			default:
				nack = true
			}

			if nack {
				if taskExecutionHistory == (model.TaskExecutionHistory{}) || job == (model.Job{}) {
					log.Info().Msg("NAcking")
					msg.Nack(false, nack)
					// continue
				}

				if taskExecutionHistory.RetryCount <= job.MaxRetryCount {
					// produce a jobId to x.retry
					taskExecutionHistory.RetryCount += 1
					err := processor.msgBroker.Publish(ctx, "x.task", "task-retry")
					if err != nil {
						log.Error().Err(err)
					}
				} else {
					// produce a Job to dlx.task
					taskExecutionHistory.Status = FAILED_STATUS
					err := processor.msgBroker.Publish(ctx, "dlx.task", "dlq.task", job)
					if err != nil {
						log.Error().Err(err)
					}
				}
				processor.DB.Save(&taskExecutionHistory)
				msg.Nack(false, nack)
			} else {
				log.Info().Msg("Acking")
				msg.Ack(false)
			}

			if cnt == 1 {
				break
			}
		}

		processor.done <- struct{}{}
	}()

	return nil
}
