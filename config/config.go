package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/exception"
)

type Config struct {
	Environment         string `mapstructure:"ENVIRONMENT"`
	DBSource            string `mapstructure:"DB_SOURCE"`
	RabbitMQServerURL   string `mapstructure:"RABBITMQ_SERVER_URL"`
	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func LoadConfig(path string) Config {
	config := Config{}
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	exception.FatalIfNeeded(err, "cannot load config")

	err = viper.Unmarshal(&config)
	exception.FatalIfNeeded(err, "failed to unmarshal config to struct")

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return config
}
