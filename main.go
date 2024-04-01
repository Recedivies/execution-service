package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/config"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/consumer"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/exception"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/mail"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/producer"
	"gorm.io/gorm"
)

func main() {
	configuration := config.LoadConfig(".")
	db := config.NewPostgres(configuration.DBSource)
	rmq := config.NewRabbitMQ(configuration.RabbitMQServerURL)
	msgBroker, err := producer.NewTask(rmq.Channel)
	exception.FatalIfNeeded(err, "rabbitmq.NewTask")

	errC, err := runTaskProcessor(configuration, db, rmq, msgBroker)
	exception.FatalIfNeeded(err, "Couldn't run")

	err = <-errC
	exception.FatalIfNeeded(err, "Error while running")
}

func runTaskProcessor(configuration config.Config, DB *gorm.DB, rmq *config.RabbitMQ, msgBroker *producer.Task) (<-chan error, error) {
	mail := mail.NewGmailSender(configuration.EmailSenderName, configuration.EmailSenderAddress, configuration.EmailSenderPassword)
	taskProcessor := consumer.NewRabbitMQTaskProcessor(mail, DB, rmq, msgBroker)

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Info().Msg("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			rmq.Close()
			stop()
			cancel()
			close(errC)
		}()

		if err := taskProcessor.Shuwdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Info().Msg("Shutdown completed")
	}()

	go func() {
		log.Info().Msg("Listening and serving")

		if err := taskProcessor.ListenAndServe(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}
