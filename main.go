package main

import (
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/api"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/config"
)

func main() {
	configuration := config.LoadConfig(".")
	database := config.NewPostgres(configuration.DBSource)

	api.RunGinServer(configuration, database)
}
